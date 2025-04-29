# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands
- **Install**: `pip install pre-commit && pre-commit install --hook-type commit-msg --hook-type post-commit`
- **Lint**: `pre-commit run --all-files`
- **Test scripts**: `bash ./claude_sanity_check.sh` or `bash ./claude_post_commit.sh`
- **Analyze logs**: `./analyze-claude-warnings.sh --list` or `./analyze-claude-warnings.sh --summary`

## Style Guidelines
- **Shell Scripts**: Use `bash` shebang, error handling with traps, strict mode (`set -eo pipefail`)
- **Formatting**: 2-space indentation, meaningful variable names with UPPER_CASE for constants
- **Error Handling**: All scripts must handle errors gracefully and provide informative messages
- **Logging**: Use structured logging format in `.claude-warnings.log`
- **Python**: Follow PEP 8 guidelines, use setuptools for packaging
- **Git Commits**: Follow conventional commits specification, all commits must pass pre-commit hooks
- **CI Awareness**: All scripts should auto-detect and skip in CI environments

## Architecture
This repo contains pre-commit hooks for Claude AI code review that enforce development philosophy and coding standards in Git workflows.