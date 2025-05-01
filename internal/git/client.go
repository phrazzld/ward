// Package git provides functionality for interacting with Git repositories
package git

import "context"

// Client defines operations for interacting with Git repositories
type Client interface {
	// GetStagedDiff retrieves the git diff of staged changes
	GetStagedDiff(ctx context.Context) (string, error)

	// GetCommitMessage retrieves the proposed commit message
	GetCommitMessage(ctx context.Context) (string, error)

	// GetCurrentBranch retrieves the name of the current git branch
	GetCurrentBranch(ctx context.Context) (string, error)

	// GetHeadCommitHash retrieves the hash of the current HEAD commit
	GetHeadCommitHash(ctx context.Context) (string, error)

	// GetGitDir retrieves the path to the .git directory
	GetGitDir(ctx context.Context) (string, error)

	// WritePlaceholderHash writes a placeholder hash to a file in the .git directory
	WritePlaceholderHash(ctx context.Context, correlationID string) error

	// ReadPlaceholderHash reads a placeholder hash from a file in the .git directory
	ReadPlaceholderHash(ctx context.Context) (string, error)

	// CleanPlaceholderHash removes the placeholder hash file from the .git directory
	CleanPlaceholderHash(ctx context.Context) error
}