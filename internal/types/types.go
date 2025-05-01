// Package types defines shared data structures used throughout the ward application
package types

import "time"

// ReviewRequest contains all information needed for code review
type ReviewRequest struct {
	// StagedDiff is the git diff of staged changes
	StagedDiff string
	// CommitMessage is the proposed commit message
	CommitMessage string
	// Branch is the current git branch
	Branch string
	// CorrelationID uniquely identifies this review request
	CorrelationID string
}

// ReviewResult represents the outcome of a code review
type ReviewResult struct {
	// Status is the review outcome: PASS, WARN, or FAIL
	Status string
	// Messages contains feedback from the review
	Messages []string
	// CorrelationID links to the original review request
	CorrelationID string
	// Timestamp when the review was completed
	Timestamp time.Time
}

// LogEntry represents a record in the ward-warnings.log file
type LogEntry struct {
	// CorrelationID uniquely identifies this review across processes
	CorrelationID string
	// Timestamp when the entry was created
	Timestamp time.Time
	// Status is the review outcome: PASS, WARN, or FAIL
	Status string
	// Messages contains feedback from the review
	Messages []string
	// CommitHash is the actual git commit hash (or placeholder)
	CommitHash string
	// IsPlaceholder indicates if CommitHash is a placeholder
	IsPlaceholder bool
}

// Config contains application configuration settings
type Config struct {
	// ClaudeCLIPath is the path to the Claude CLI executable
	ClaudeCLIPath string
	// LogFilePath is the path to the ward-warnings.log file
	LogFilePath string
	// PromptTemplate for LLM code review
	PromptTemplate string
	// SkipInCI indicates whether to skip reviews in CI environments
	SkipInCI bool
}
