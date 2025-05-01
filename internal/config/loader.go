// Package config provides functionality for loading and validating application configuration
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	werrors "github.com/phrazzld/ward/internal/errors"
	"github.com/phrazzld/ward/internal/types"
)

const (
	// Environment variable names
	envClaudeCLIPath   = "WARD_CLAUDE_CLI_PATH"
	envLogFilePath     = "WARD_LOG_FILE_PATH"
	envPromptTemplate  = "WARD_PROMPT_TEMPLATE"
	envSkipInCI        = "WARD_SKIP_IN_CI"
)

// Load reads configuration from environment variables and validates required fields
func Load() (*types.Config, error) {
	config := &types.Config{}

	// Load and validate ClaudeCLIPath (required)
	claudeCLIPath := os.Getenv(envClaudeCLIPath)
	if claudeCLIPath == "" {
		return nil, werrors.WrapConfigError(
			errors.New("missing required value"),
			"ClaudeCLIPath",
		)
	}
	config.ClaudeCLIPath = claudeCLIPath

	// Load and validate LogFilePath (required)
	logFilePath := os.Getenv(envLogFilePath)
	if logFilePath == "" {
		return nil, werrors.WrapConfigError(
			errors.New("missing required value"),
			"LogFilePath",
		)
	}
	config.LogFilePath = logFilePath

	// Load and validate PromptTemplate (required)
	promptTemplate := os.Getenv(envPromptTemplate)
	if promptTemplate == "" {
		return nil, werrors.WrapConfigError(
			errors.New("missing required value"),
			"PromptTemplate",
		)
	}
	config.PromptTemplate = promptTemplate

	// Load and parse SkipInCI (optional, defaults to false)
	skipInCIStr := os.Getenv(envSkipInCI)
	if skipInCIStr != "" {
		skipInCI, err := parseBool(skipInCIStr)
		if err != nil {
			return nil, werrors.WrapConfigError(
				fmt.Errorf("invalid boolean value '%s': %w", skipInCIStr, err),
				"SkipInCI",
			)
		}
		config.SkipInCI = skipInCI
	}

	return config, nil
}

// parseBool parses a string into a boolean, accepting various representations
func parseBool(str string) (bool, error) {
	// First try the standard Go parsing
	val, err := strconv.ParseBool(str)
	if err == nil {
		return val, nil
	}

	// Handle additional formats
	switch strings.ToLower(strings.TrimSpace(str)) {
	case "yes", "y", "1":
		return true, nil
	case "no", "n", "0":
		return false, nil
	default:
		return false, fmt.Errorf("cannot parse '%s' as boolean", str)
	}
}
