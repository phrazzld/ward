// Package util provides utility functions used throughout the ward application
package util

import (
	"context"
	"crypto/rand"
	"fmt"
	"os/exec"

	"github.com/phrazzld/ward/internal/errors"
)

// GenerateCorrelationID creates a unique UUID v4-like identifier for request tracking
// The returned string is a hex-encoded random value suitable for correlation IDs
func GenerateCorrelationID() (string, error) {
	// Generate 16 random bytes (UUID v4 size)
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Set version (4) and variant bits according to RFC 4122
	bytes[6] = (bytes[6] & 0x0F) | 0x40 // Version 4
	bytes[8] = (bytes[8] & 0x3F) | 0x80 // Variant 1

	// Format as standard UUID string
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:16])

	return uuid, nil
}

// ExecCommand is a safe wrapper around exec.CommandContext that avoids shell interpolation
// risks and provides consistent error handling
// ctx: Context for command execution, can be used for timeouts/cancellation
// command: The executable to run
// args: Arguments to pass to the command (passed directly, no shell expansion)
func ExecCommand(ctx context.Context, command string, args ...string) ([]byte, error) {
	// Create command with proper context
	cmd := exec.CommandContext(ctx, command, args...)

	// Execute and capture output (combines stdout and stderr)
	output, err := cmd.CombinedOutput()

	// Handle errors with consistent formatting using sentinel error
	if err != nil {
		return output, errors.WrapExternalCmdError(err, command, args)
	}

	return output, nil
}
