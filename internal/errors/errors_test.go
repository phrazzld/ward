package errors_test

import (
	"errors"
	"testing"

	werrors "github.com/phrazzld/ward/internal/errors"
)

func TestErrorsAreAddressable(t *testing.T) {
	// Verify that errors can be used from another package
	var err error

	// Test ErrCheckFailed
	err = werrors.ErrCheckFailed
	if !errors.Is(err, werrors.ErrCheckFailed) {
		t.Errorf("Expected error to be ErrCheckFailed")
	}

	// Test ErrLogUpdateNotFound
	err = werrors.ErrLogUpdateNotFound
	if !errors.Is(err, werrors.ErrLogUpdateNotFound) {
		t.Errorf("Expected error to be ErrLogUpdateNotFound")
	}

	// Test ErrExternalCmdFailed
	err = werrors.ErrExternalCmdFailed
	if !errors.Is(err, werrors.ErrExternalCmdFailed) {
		t.Errorf("Expected error to be ErrExternalCmdFailed")
	}

	// Test ErrConfigInvalid
	err = werrors.ErrConfigInvalid
	if !errors.Is(err, werrors.ErrConfigInvalid) {
		t.Errorf("Expected error to be ErrConfigInvalid")
	}

	// Test error wrapping
	testErr := errors.New("test error")
	wrappedErr := werrors.WrapExternalCmdError(testErr, "test", []string{"arg1", "arg2"})
	if !errors.Is(wrappedErr, werrors.ErrExternalCmdFailed) {
		t.Errorf("Wrapped error should be ErrExternalCmdFailed")
	}

	wrappedConfigErr := werrors.WrapConfigError(testErr, "testField")
	if !errors.Is(wrappedConfigErr, werrors.ErrConfigInvalid) {
		t.Errorf("Wrapped config error should be ErrConfigInvalid")
	}
}
