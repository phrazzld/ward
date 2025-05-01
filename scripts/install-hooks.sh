#!/usr/bin/env bash

# Install pre-commit hooks script
set -eo pipefail

# Detect if script is being run in CI environment
if [ -n "${CI:-}" ] || [ -n "${GITHUB_ACTIONS:-}" ] || [ -n "${GITLAB_CI:-}" ] || [ -n "${CIRCLECI:-}" ]; then
  echo "CI environment detected, skipping pre-commit hook installation"
  exit 0
fi

echo "==> Installing pre-commit hooks..."

# First check if pre-commit is installed
if ! command -v pre-commit &> /dev/null; then
  echo "pre-commit is not installed, attempting to install..."
  if command -v pip &> /dev/null; then
    pip install pre-commit
  elif command -v pip3 &> /dev/null; then
    pip3 install pre-commit
  else
    echo "ERROR: pip not found. Please install pip and pre-commit manually."
    echo "Run: pip install pre-commit"
    exit 1
  fi
fi

# Check if gocyclo is installed, install if not
if ! command -v gocyclo &> /dev/null; then
  echo "gocyclo is not installed, attempting to install..."
  if command -v go &> /dev/null; then
    go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
  else
    echo "ERROR: go not found. Please install go and gocyclo manually."
    exit 1
  fi
fi

# Install the pre-commit hooks
pre-commit install --hook-type pre-commit --hook-type commit-msg --hook-type post-commit

# Create or update the ward-specific pre-commit hook chain
HOOK_PATH=".git/hooks/post-commit"

# Ensure the hook exists and is executable
if [ ! -x "$HOOK_PATH" ]; then
  echo "Creating post-commit hook with ward-log..."
else
  # Check if the ward-log is already in the hook
  if grep -q "ward_log.sh" "$HOOK_PATH"; then
    echo "ward_log.sh already in post-commit hook"
  else
    # Backup the current hook
    cp "$HOOK_PATH" "$HOOK_PATH.backup"

    # Add ward_log.sh to the hook
    echo -e "#!/usr/bin/env bash\n\n# Run original pre-commit generated hook\n$(cat $HOOK_PATH)\n\n# Run Ward log hook\nbash ./ward_log.sh" > "$HOOK_PATH"
    chmod +x "$HOOK_PATH"
  fi
fi

echo "==> Hooks installed successfully!"
