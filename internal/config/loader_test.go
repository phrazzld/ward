package config

import (
	"errors"
	"testing"

	werrors "github.com/phrazzld/ward/internal/errors"
)

func TestLoadSuccess(t *testing.T) {
	t.Run("all required fields set", func(t *testing.T) {
		t.Setenv(envClaudeCLIPath, "/path/to/claude")
		t.Setenv(envLogFilePath, "/path/to/log")
		t.Setenv(envPromptTemplate, "template")

		config, err := Load()
		if err != nil {
			t.Fatalf("Load() returned unexpected error: %v", err)
		}

		if config.ClaudeCLIPath != "/path/to/claude" {
			t.Errorf("expected ClaudeCLIPath to be '/path/to/claude', got '%s'", config.ClaudeCLIPath)
		}
		if config.LogFilePath != "/path/to/log" {
			t.Errorf("expected LogFilePath to be '/path/to/log', got '%s'", config.LogFilePath)
		}
		if config.PromptTemplate != "template" {
			t.Errorf("expected PromptTemplate to be 'template', got '%s'", config.PromptTemplate)
		}
		if config.SkipInCI != false {
			t.Errorf("expected SkipInCI to default to false, got %v", config.SkipInCI)
		}
	})

	t.Run("all fields including optional set", func(t *testing.T) {
		t.Setenv(envClaudeCLIPath, "/path/to/claude")
		t.Setenv(envLogFilePath, "/path/to/log")
		t.Setenv(envPromptTemplate, "template")
		t.Setenv(envSkipInCI, "true")

		config, err := Load()
		if err != nil {
			t.Fatalf("Load() returned unexpected error: %v", err)
		}

		if !config.SkipInCI {
			t.Errorf("expected SkipInCI to be true, got %v", config.SkipInCI)
		}
	})
}

func TestLoadErrorMissingRequired(t *testing.T) {
	testCases := []struct {
		name           string
		claudeCLIPath  string
		logFilePath    string
		promptTemplate string
		expectedField  string
	}{
		{
			name:           "missing ClaudeCLIPath",
			claudeCLIPath:  "",
			logFilePath:    "/path/to/log",
			promptTemplate: "template",
			expectedField:  "ClaudeCLIPath",
		},
		{
			name:           "missing LogFilePath",
			claudeCLIPath:  "/path/to/claude",
			logFilePath:    "",
			promptTemplate: "template",
			expectedField:  "LogFilePath",
		},
		{
			name:           "missing PromptTemplate",
			claudeCLIPath:  "/path/to/claude",
			logFilePath:    "/path/to/log",
			promptTemplate: "",
			expectedField:  "PromptTemplate",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(envClaudeCLIPath, tc.claudeCLIPath)
			t.Setenv(envLogFilePath, tc.logFilePath)
			t.Setenv(envPromptTemplate, tc.promptTemplate)

			_, err := Load()
			if err == nil {
				t.Fatal("Load() expected to return error, got nil")
			}

			if !errors.Is(err, werrors.ErrConfigInvalid) {
				t.Errorf("expected error to be ErrConfigInvalid, got %v", err)
			}

			expectedMsg := "configuration is invalid: field '" + tc.expectedField + "': missing required value"
			if err.Error() != expectedMsg {
				t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
			}
		})
	}
}

func TestLoadErrorInvalidOptional(t *testing.T) {
	t.Setenv(envClaudeCLIPath, "/path/to/claude")
	t.Setenv(envLogFilePath, "/path/to/log")
	t.Setenv(envPromptTemplate, "template")
	t.Setenv(envSkipInCI, "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected to return error, got nil")
	}

	if !errors.Is(err, werrors.ErrConfigInvalid) {
		t.Errorf("expected error to be ErrConfigInvalid, got %v", err)
	}
}

func TestParseBool(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
		wantErr  bool
	}{
		// Standard boolean strings
		{"true", true, false},
		{"false", false, false},
		{"t", true, false},
		{"f", false, false},
		{"TRUE", true, false},
		{"FALSE", false, false},
		{"True", true, false},
		{"False", false, false},

		// Additional formats
		{"yes", true, false},
		{"no", false, false},
		{"y", true, false},
		{"n", false, false},
		{"YES", true, false},
		{"NO", false, false},
		{"Y", true, false},
		{"N", false, false},
		{"1", true, false},
		{"0", false, false},
		{" yes ", true, false}, // Test with whitespace

		// Invalid formats
		{"invalid", false, true},
		{"truee", false, true},
		{"2", false, true},
		{"", false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result, err := parseBool(tc.input)

			// Check error case
			if tc.wantErr && err == nil {
				t.Errorf("parseBool(%q) expected error, got nil", tc.input)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("parseBool(%q) unexpected error: %v", tc.input, err)
			}

			// Check result value
			if !tc.wantErr && result != tc.expected {
				t.Errorf("parseBool(%q) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
