# Ward - Code Review Hooks

Pre-commit hooks for integrating code review into your Git workflow.

## Features

- Automated code review before each commit
- Logs warnings and failures with correlation IDs
- Post-commit updating of log entries with actual commit hashes
- Analysis tool for reviewing warnings and failures

## Installation

1. Install the `pre-commit` framework:
   ```
   pip install pre-commit
   ```

2. Add to your `.pre-commit-config.yaml`:
   ```yaml
   repos:
   - repo: https://github.com/phrazzld/ward
     rev: v0.1.0  # Use the latest tag
     hooks:
     - id: ward-check
     - id: ward-log
   ```

3. Install the hooks:
   ```
   pre-commit install --hook-type commit-msg --hook-type post-commit
   ```

## Requirements

- Claude CLI must be installed and configured (https://github.com/anthropics/claude-cli)

## Usage

### Automatic Usage

Once installed, the hooks will run automatically:
- `ward-check`: Runs before each commit is created
- `ward-log`: Runs after each commit to update log entries

### Log Analysis

Use the included analysis script to review warnings and failures:

```
./ward_analyze.sh --help
```

Options:
- `--list`: List all entries
- `-c HASH`: Show details for a specific commit
- `-s STATUS`: Filter by status (WARN, FAIL)
- `-b BRANCH`: Filter by branch
- `--summary`: Show summary statistics

## Configuration

The hooks automatically skip in CI environments. No additional configuration is needed.

## Log File

Warnings and failures are logged to `.ward-warnings.log` in your project root.
Add this file to your `.gitignore` to prevent accidental commits.

## License

MIT