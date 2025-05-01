// Package errors defines sentinel errors used throughout the ward application
package errors

import (
	"errors"
	"fmt"
)

// Standard errors that can be used for equality checks
var (
	// ErrCheckFailed is returned when a code review check fails
	ErrCheckFailed = errors.New("code review check failed")

	// ErrLogUpdateNotFound is returned when trying to update a log entry that doesn't exist
	ErrLogUpdateNotFound = errors.New("log entry not found for update")

	// ErrExternalCmdFailed is returned when an external command execution fails
	ErrExternalCmdFailed = errors.New("external command execution failed")

	// ErrConfigInvalid is returned when configuration is invalid
	ErrConfigInvalid = errors.New("configuration is invalid")
)

// WrapExternalCmdError adds context to a command error
func WrapExternalCmdError(err error, cmd string, args []string) error {
	return fmt.Errorf("%w: '%s %v': %w", ErrExternalCmdFailed, cmd, args, err)
}

// WrapConfigError adds context to a configuration error
func WrapConfigError(err error, field string) error {
	return fmt.Errorf("%w: field '%s': %w", ErrConfigInvalid, field, err)
}
