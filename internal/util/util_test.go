package util

import (
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	werrors "github.com/phrazzld/ward/internal/errors"
)

func TestGenerateCorrelationID(t *testing.T) {
	// Define the expected format of a UUID v4
	// Format: 8-4-4-4-12 characters (hex), where certain positions have specific values
	uuidv4Regex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

	t.Run("generates valid UUID v4 format", func(t *testing.T) {
		id, err := GenerateCorrelationID()
		if err != nil {
			t.Fatalf("GenerateCorrelationID() error = %v", err)
		}

		if !uuidv4Regex.MatchString(id) {
			t.Errorf("GenerateCorrelationID() = %v, doesn't match UUID v4 format", id)
		}

		// Check version bits (6th character of 3rd group should be 4)
		segments := strings.Split(id, "-")
		if len(segments) != 5 {
			t.Errorf("UUID doesn't have 5 segments: %v", id)
		} else if !strings.HasPrefix(segments[2], "4") {
			t.Errorf("UUID version bits not set correctly: %v", id)
		}

		// Check variant bits (first character of 4th group should be 8, 9, a, or b)
		if len(segments) == 5 {
			firstChar := segments[3][0:1]
			if firstChar != "8" && firstChar != "9" && firstChar != "a" && firstChar != "b" {
				t.Errorf("UUID variant bits not set correctly: %v", id)
			}
		}
	})

	t.Run("generates unique IDs", func(t *testing.T) {
		// Generate multiple IDs and ensure they're unique
		ids := make(map[string]bool)
		iterations := 5 // Reasonable number for a unit test

		for i := 0; i < iterations; i++ {
			id, err := GenerateCorrelationID()
			if err != nil {
				t.Fatalf("GenerateCorrelationID() error = %v", err)
			}

			if ids[id] {
				t.Errorf("Generated duplicate ID: %v", id)
			}
			ids[id] = true
		}
	})
}

// testSuccessfulExecution tests a successful command execution
func testSuccessfulExecution(t *testing.T) {
	ctx := context.Background()
	output, err := ExecCommand(ctx, "echo", "test")
	if err != nil {
		t.Fatalf("ExecCommand() error = %v", err)
	}

	expected := "test\n" // echo adds a newline
	if string(output) != expected {
		t.Errorf("ExecCommand() = %q, want %q", output, expected)
	}
}

// testErrorExitCode tests a command that exits with an error code
func testErrorExitCode(t *testing.T) {
	ctx := context.Background()
	// Use a command that will fail with a non-zero exit code
	// 'false' is a simple command that always exits with code 1
	output, err := ExecCommand(ctx, "false")

	// Check that error is not nil
	if err == nil {
		t.Fatalf("ExecCommand() expected error, got nil")
	}

	// Verify error is wrapped correctly
	if !errors.Is(err, werrors.ErrExternalCmdFailed) {
		t.Errorf("ExecCommand() error = %v, should wrap ErrExternalCmdFailed", err)
	}

	// Output should be empty for false
	if len(output) > 0 {
		t.Errorf("ExecCommand() output = %q, want empty", output)
	}
}

// testStderrOutput tests a command that generates stderr output
func testStderrOutput(t *testing.T) {
	ctx := context.Background()
	// The command 'ls /nonexistent' will generate an error message to stderr
	output, err := ExecCommand(ctx, "ls", "/nonexistent")

	// Check that error is not nil
	if err == nil {
		t.Fatalf("ExecCommand() expected error, got nil")
	}

	// Verify error is wrapped correctly
	if !errors.Is(err, werrors.ErrExternalCmdFailed) {
		t.Errorf("ExecCommand() error = %v, should wrap ErrExternalCmdFailed", err)
	}

	// Output should contain error message from stderr
	if !strings.Contains(string(output), "No such file or directory") {
		t.Errorf("ExecCommand() output = %q, should contain stderr message", output)
	}
}

// testCommandTimeout tests a command that times out
func testCommandTimeout(t *testing.T) {
	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Use a command that will take longer than the timeout
	// 'sleep' is available on most systems
	_, err := ExecCommand(ctx, "sleep", "1")

	// Check that error is not nil
	if err == nil {
		t.Fatalf("ExecCommand() expected error, got nil")
	}

	// Verify error is wrapped correctly
	if !errors.Is(err, werrors.ErrExternalCmdFailed) {
		t.Errorf("ExecCommand() error = %v, should wrap ErrExternalCmdFailed", err)
	}

	// Verify we got a context deadline exceeded error inside
	if !errors.Is(err, context.DeadlineExceeded) && !strings.Contains(err.Error(), "context deadline exceeded") {
		// On some systems, we might not get a direct context error but a process killed message
		if !strings.Contains(err.Error(), "signal: killed") && !strings.Contains(err.Error(), "process killed") {
			t.Errorf("ExecCommand() error = %v, should be context deadline or process killed", err)
		}
	}

	// We don't capture output for timeout tests as it's unpredictable
}

// testNonExistentCommand tests a command that doesn't exist
func testNonExistentCommand(t *testing.T) {
	ctx := context.Background()
	// Use a command that doesn't exist
	output, err := ExecCommand(ctx, "command_that_definitely_does_not_exist")

	// Check that error is not nil
	if err == nil {
		t.Fatalf("ExecCommand() expected error, got nil")
	}

	// Verify error is wrapped correctly
	if !errors.Is(err, werrors.ErrExternalCmdFailed) {
		t.Errorf("ExecCommand() error = %v, should wrap ErrExternalCmdFailed", err)
	}

	// Verify the specific executable not found error
	var execErr *exec.Error
	if !errors.As(err, &execErr) && !strings.Contains(err.Error(), "executable file not found") {
		t.Errorf("ExecCommand() should have exec.Error or 'not found' message for non-existent command")
	}

	// Output should be empty
	if len(output) > 0 {
		t.Errorf("ExecCommand() output = %q, want empty", output)
	}
}

func TestExecCommand(t *testing.T) {
	t.Run("successful command execution", testSuccessfulExecution)
	t.Run("command with error exit code", testErrorExitCode)
	t.Run("command with stderr output", testStderrOutput)
	t.Run("command with timeout", testCommandTimeout)
	t.Run("non-existent command", testNonExistentCommand)
}
