// Package shim provides functionality for generating Git hook scripts
package shim

import "context"

// Generator defines operations for creating Git hook scripts
type Generator interface {
	// Generate creates a hook script for the specified manager and hook type
	// manager is the Git hook management system (e.g., "pre-commit")
	// hookType is the type of Git hook (e.g., "pre-commit", "post-commit")
	Generate(ctx context.Context, manager string, hookType string) (string, error)
}