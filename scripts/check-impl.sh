#!/usr/bin/env bash
# Simple script to verify our implementation

set -eo pipefail

# First, run the linter on util package
echo "Running linter on util package..."
cd "$(git rev-parse --show-toplevel)"
golangci-lint run ./internal/util/...

# Check if the code builds
echo "Building util package..."
go build ./internal/util/...

# Print success message
echo "✅ Util package implementation passes checks!"
echo "The following features are now available:"
echo "- GenerateCorrelationID(): Creates RFC 4122 UUID v4-like identifiers"
echo "- ExecCommand(): Safe wrapper for external command execution with proper error handling"

# Note for next steps
echo
echo "Ready for the next task: T017 · Test · P2: add unit tests for util"
