// Package ci provides functionality for detecting continuous integration environments
package ci

import (
	"os"
)

// DefaultDetector is the default implementation of the Detector interface
type DefaultDetector struct{}

// NewDetector creates a new instance of DefaultDetector
func NewDetector() *DefaultDetector {
	return &DefaultDetector{}
}

// IsCI returns true if the code is running in a continuous integration environment,
// false otherwise. It detects CI environments by checking for common CI environment
// variables like CI and GITHUB_ACTIONS.
func (d *DefaultDetector) IsCI() bool {
	// Check for generic CI environment variable (set by most CI systems)
	if _, exists := os.LookupEnv("CI"); exists {
		return true
	}

	// Check for GitHub Actions specifically
	if _, exists := os.LookupEnv("GITHUB_ACTIONS"); exists {
		return true
	}

	// Return false if no CI environment variables were detected
	return false
}
