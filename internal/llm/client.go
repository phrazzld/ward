// Package llm provides functionality for interacting with language models
package llm

import (
	"context"

	"github.com/phrazzld/ward/internal/types"
)

// Client defines operations for interacting with language models for code review
type Client interface {
	// AnalyzeChanges sends code changes to an LLM for review and returns analysis results
	AnalyzeChanges(ctx context.Context, req *types.ReviewRequest) (*types.ReviewResult, error)
}