//go:build integration
// +build integration

package commands

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jongio/azd-core/testutil"
	"github.com/jongio/azd-exec/cli/src/internal/version"
)

func TestVersionCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name       string
		outputFlag string
		wantText   string
	}{
		{
			name:       "Default output",
			outputFlag: "default",
			wantText:   version.Version,
		},
		{
			name:       "JSON output",
			outputFlag: "json",
			wantText:   `"version"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFormat := tt.outputFlag
			cmd := NewVersionCommand(&outputFormat)

			// Capture output
			output := testutil.CaptureOutput(t, func() error {
				return cmd.Execute()
			})

			// Verify output
			if !strings.Contains(output, tt.wantText) {
				t.Errorf("Output does not contain expected text.\nGot: %s\nWant substring: %s", output, tt.wantText)
			}

			// For JSON output, verify it's valid JSON
			if tt.outputFlag == "json" {
				var result map[string]string
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("JSON output is invalid: %v\nOutput: %s", err, output)
				}
				if result["version"] == "" {
					t.Error("JSON output missing version field")
				}
			}
		})
	}
}

func TestVersionCommandIntegration_DefaultFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)

	output := testutil.CaptureOutput(t, func() error {
		return cmd.Execute()
	})

	// Default output should contain the version
	output = strings.TrimSpace(output)
	if !strings.Contains(output, version.Version) {
		t.Errorf("Default output should contain version, got: %s", output)
	}
}

func TestVersionCommandIntegration_JSONFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	outputFormat := "json"
	cmd := NewVersionCommand(&outputFormat)

	output := testutil.CaptureOutput(t, func() error {
		return cmd.Execute()
	})

	// Parse JSON
	var result map[string]string
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, output)
	}

	// Verify version field exists
	version, ok := result["version"]
	if !ok {
		t.Error("JSON output missing 'version' field")
	}

	// Verify version format (should be semver-like)
	if !strings.Contains(version, ".") {
		t.Errorf("Version should contain dots (semver format), got: %s", version)
	}
}
