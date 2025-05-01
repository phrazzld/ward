package ci

import (
	"os"
	"testing"
)

func TestIsCI(t *testing.T) {
	// Helper function to handle env var operations
	setEnv := func(key, value string) {
		if err := os.Setenv(key, value); err != nil {
			t.Fatalf("Failed to set environment variable %s: %v", key, err)
		}
	}

	unsetEnv := func(key string) {
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("Failed to unset environment variable %s: %v", key, err)
		}
	}

	// Create a table of test cases
	tests := []struct {
		name       string
		setupEnv   func()
		cleanupEnv func()
		want       bool
	}{
		{
			name: "no CI env vars set",
			setupEnv: func() {
				unsetEnv("CI")
				unsetEnv("GITHUB_ACTIONS")
			},
			cleanupEnv: func() {},
			want:       false,
		},
		{
			name: "CI env var set",
			setupEnv: func() {
				unsetEnv("GITHUB_ACTIONS")
				setEnv("CI", "true")
			},
			cleanupEnv: func() {
				unsetEnv("CI")
			},
			want: true,
		},
		{
			name: "GITHUB_ACTIONS env var set",
			setupEnv: func() {
				unsetEnv("CI")
				setEnv("GITHUB_ACTIONS", "true")
			},
			cleanupEnv: func() {
				unsetEnv("GITHUB_ACTIONS")
			},
			want: true,
		},
		{
			name: "both CI and GITHUB_ACTIONS env vars set",
			setupEnv: func() {
				setEnv("CI", "true")
				setEnv("GITHUB_ACTIONS", "true")
			},
			cleanupEnv: func() {
				unsetEnv("CI")
				unsetEnv("GITHUB_ACTIONS")
			},
			want: true,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			tt.setupEnv()
			defer tt.cleanupEnv()

			// Create detector and test
			detector := NewDetector()
			if got := detector.IsCI(); got != tt.want {
				t.Errorf("DefaultDetector.IsCI() = %v, want %v", got, tt.want)
			}
		})
	}
}
