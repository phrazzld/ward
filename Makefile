# Makefile for Ward - Code Review Hooks

# Variables
BINARY_NAME := ward
BUILD_DIR := ./bin
CMD_DIR := ./cmd/ward

# Go build flags
GO_BUILD_FLAGS := -v

# Define the default target when just `make` is executed
.DEFAULT_GOAL := help

# Help message
.PHONY: help
help:
	@echo "Ward Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build    Build the ward binary"
	@echo "  make test     Run tests across all packages"
	@echo "  make lint     Run linting checks"
	@echo "  make clean    Remove build artifacts"
	@echo "  make help     Show this help message"
	@echo ""

# Build the binary
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build $(GO_BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

# Run tests
.PHONY: test
test:
	go test ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# All - Run lint, test, and build
.PHONY: all
all: lint test build