// Package git provides functionality for interacting with Git repositories
package git

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/phrazzld/ward/internal/errors"
	"github.com/phrazzld/ward/internal/util"
)

// GitClient is an implementation of the Client interface that uses
// git command-line interface via the util.ExecCommand wrapper
type GitClient struct{}

// NewClient creates a new GitClient instance
func NewClient() *GitClient {
	return &GitClient{}
}

// GetStagedDiff retrieves the git diff of staged changes
func (c *GitClient) GetStagedDiff(ctx context.Context) (string, error) {
	output, err := util.ExecCommand(ctx, "git", "diff", "--cached")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GetCommitMessage retrieves the proposed commit message
func (c *GitClient) GetCommitMessage(ctx context.Context) (string, error) {
	// Get the git directory
	gitDir, err := c.GetGitDir(ctx)
	if err != nil {
		return "", err
	}

	// Read COMMIT_EDITMSG file
	commitMsgPath := filepath.Join(gitDir, "COMMIT_EDITMSG")
	data, err := os.ReadFile(commitMsgPath)
	if err != nil {
		return "", fmt.Errorf("failed to read commit message from %s: %w", commitMsgPath, err)
	}

	return string(data), nil
}

// GetCurrentBranch retrieves the name of the current git branch
func (c *GitClient) GetCurrentBranch(ctx context.Context) (string, error) {
	output, err := util.ExecCommand(ctx, "git", "symbolic-ref", "--short", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetHeadCommitHash retrieves the hash of the current HEAD commit
func (c *GitClient) GetHeadCommitHash(ctx context.Context) (string, error) {
	output, err := util.ExecCommand(ctx, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitDir retrieves the path to the .git directory
func (c *GitClient) GetGitDir(ctx context.Context) (string, error) {
	output, err := util.ExecCommand(ctx, "git", "rev-parse", "--git-dir")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// getPlaceholderPath returns the path to the placeholder file
func (c *GitClient) getPlaceholderPath(ctx context.Context) (string, error) {
	gitDir, err := c.GetGitDir(ctx)
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, "ward_placeholder"), nil
}

// WritePlaceholderHash writes a placeholder hash to a file in the .git directory
func (c *GitClient) WritePlaceholderHash(ctx context.Context, correlationID string) error {
	placeholderPath, err := c.getPlaceholderPath(ctx)
	if err != nil {
		return err
	}

	// Write correlation ID to the placeholder file
	err = os.WriteFile(placeholderPath, []byte(correlationID), 0644)
	if err != nil {
		return fmt.Errorf("failed to write placeholder hash: %w", err)
	}

	return nil
}

// ReadPlaceholderHash reads a placeholder hash from a file in the .git directory
func (c *GitClient) ReadPlaceholderHash(ctx context.Context) (string, error) {
	placeholderPath, err := c.getPlaceholderPath(ctx)
	if err != nil {
		return "", err
	}

	// Check if placeholder file exists
	if _, err := os.Stat(placeholderPath); os.IsNotExist(err) {
		return "", fmt.Errorf("%w: placeholder file not found", errors.ErrLogUpdateNotFound)
	}

	// Read correlation ID from the placeholder file
	data, err := os.ReadFile(placeholderPath)
	if err != nil {
		return "", fmt.Errorf("failed to read placeholder hash: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

// CleanPlaceholderHash removes the placeholder hash file from the .git directory
func (c *GitClient) CleanPlaceholderHash(ctx context.Context) error {
	placeholderPath, err := c.getPlaceholderPath(ctx)
	if err != nil {
		return err
	}

	// Check if file exists before attempting to remove
	if _, err := os.Stat(placeholderPath); os.IsNotExist(err) {
		// File doesn't exist, so nothing to clean
		return nil
	}

	// Remove the placeholder file
	err = os.Remove(placeholderPath)
	if err != nil {
		return fmt.Errorf("failed to remove placeholder hash file: %w", err)
	}

	return nil
}
