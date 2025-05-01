// Package log provides functionality for writing and managing log entries
package log

import (
	"context"

	"github.com/phrazzld/ward/internal/types"
)

// Writer defines operations for writing and updating log entries
type Writer interface {
	// WriteEntry writes a new log entry to the ward warnings log file
	WriteEntry(ctx context.Context, entry *types.LogEntry) error

	// UpdateEntryCommitHash updates the commit hash of an existing log entry
	// identified by its correlation ID
	UpdateEntryCommitHash(ctx context.Context, correlationID, commitHash string) error
}
