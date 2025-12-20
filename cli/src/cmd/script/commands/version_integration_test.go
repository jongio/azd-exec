//go:build integration
// +build integration

package commands

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/jongio/azd-exec/cli/src/internal/testhelpers"
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
			wantText:   "azd exec version",
		},
		{
			name:       "JSON output",
			outputFlag: "json",
			wantText:   `"version"`,
		},
		{
			name:       "Quiet output",
			outputFlag: "",
			wantText:   "0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFormat := tt.outputFlag
			cmd := NewVersionCommand(&outputFormat)

			// Set quiet flag if needed
			if tt.name == "Quiet output" {
				cmd.Flags().Set("quiet", "true")
			}

			// Capture output
			output := testhelpers.CaptureOutput(t, func() error {
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

func TestVersionCommandIntegration_QuietFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)
	cmd.Flags().Set("quiet", "true")

	output := testhelpers.CaptureOutput(t, func() error {
		return cmd.Execute()
	})

	// Quiet output should be just the version number
	output = strings.TrimSpace(output)
	if !strings.HasPrefix(output, "0.") {
		t.Errorf("Quiet output should be just version number, got: %s", output)
	}

	// Should not contain "azd exec version" prefix
	if strings.Contains(output, "azd exec version") {
		t.Errorf("Quiet output should not contain prefix, got: %s", output)
	}
}

func TestVersionCommandIntegration_JSONFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	outputFormat := "json"
	cmd := NewVersionCommand(&outputFormat)

	output := testhelpers.CaptureOutput(t, func() error {
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
