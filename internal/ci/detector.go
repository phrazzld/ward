// Package ci provides functionality for detecting continuous integration environments
package ci

// Detector defines operations for detecting if code is running in a CI environment
type Detector interface {
	// IsCI returns true if the code is running in a continuous integration environment,
	// false otherwise
	IsCI() bool
}