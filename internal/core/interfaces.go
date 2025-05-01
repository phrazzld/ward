// Package core provides the central interfaces for the ward application
package core

import (
	"context"
)

// Checker defines operations for checking code changes and providing review feedback
type Checker interface {
	// Check examines staged changes and provides review feedback
	// It retrieves staged changes, sends them to LLM for analysis,
	// logs the results, and manages the placeholder hash
	Check(ctx context.Context) error
}

// LogUpdater defines operations for updating log entries with commit hashes
type LogUpdater interface {
	// Update reads the placeholder hash, gets the actual commit hash,
	// updates the corresponding log entry, and cleans up the placeholder
	Update(ctx context.Context) error
}

// Initializer defines operations for setting up Git hooks
type Initializer interface {
	// Init generates and installs Git hook scripts for the specified hook manager
	// hookManager is the Git hook management system (e.g., "pre-commit")
	Init(ctx context.Context, hookManager string) error
}
