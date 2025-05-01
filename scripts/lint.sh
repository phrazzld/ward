#!/usr/bin/env bash
# Run golangci-lint with our desired configurations directly via CLI arguments

set -eo pipefail

# Error handling function
handle_error() {
  echo "Error occurred in lint script at line $1"
  exit 1
}

# Set up trap to catch errors
trap 'handle_error $LINENO' ERR

# Create directory for the script
mkdir -p "$(dirname "$0")"

echo "Running Go linters..."

# Run golangci-lint with explicitly enabled linters
golangci-lint run \
  --enable=errcheck \
  --enable=govet \
  --enable=staticcheck \
  --enable=unused \
  --enable=ineffassign \
  --enable=gocritic \
  --enable=stylecheck \
  --enable=gofmt \
  --enable=goimports \
  --enable=nolintlint \
  --enable=gocyclo \
  --enable=bodyclose \
  --enable=exportloopref \
  --enable=unconvert \
  --enable=unparam \
  "$@"

exit $?
