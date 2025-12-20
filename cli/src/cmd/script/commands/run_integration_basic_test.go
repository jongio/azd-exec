//go:build integration
// +build integration

package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jongio/azd-exec/cli/src/internal/testhelpers"
)

func TestRunCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get path to test projects
	testProjectsDir := testhelpers.GetTestProjectsDir(t)

	tests := []struct {
		name        string
		scriptPath  string
		args        []string
		skipWindows bool
		skipUnix    bool
		wantOutput  string
	}{
		{
			name:        "Bash simple script",
			scriptPath:  filepath.Join(testProjectsDir, "bash", "simple.sh"),
			skipWindows: true, // Skip bash tests on Windows
			wantOutput:  "Hello from bash script",
		},
		{
			name:        "Bash script with args",
			scriptPath:  filepath.Join(testProjectsDir, "bash", "with-args.sh"),
			args:        []string{"test1", "test2"},
			skipWindows: true,
			wantOutput:  "Arg 1: test1",
		},
		{
			name:        "Bash env test",
			scriptPath:  filepath.Join(testProjectsDir, "bash", "env-test.sh"),
			skipWindows: true,
			wantOutput:  "PATH exists: yes",
		},
		{
			name:       "PowerShell simple script",
			scriptPath: filepath.Join(testProjectsDir, "powershell", "simple.ps1"),
			wantOutput: "Hello from PowerShell script",
		},
		{
			name:       "PowerShell env test",
			scriptPath: filepath.Join(testProjectsDir, "powershell", "env-test.ps1"),
			wantOutput: "PATH exists: True",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipWindows && runtime.GOOS == "windows" {
				t.Skip("Skipping test on Windows")
			}
			if tt.skipUnix && runtime.GOOS != "windows" {
				t.Skip("Skipping test on Unix")
			}

			// Verify script exists
			if _, err := os.Stat(tt.scriptPath); os.IsNotExist(err) {
				t.Fatalf("Test script not found: %s", tt.scriptPath)
			}

			// Create command
			outputFormat := "default"
			cmd := NewRunCommand(&outputFormat)

			// Set args
			cmdArgs := []string{tt.scriptPath}
			cmdArgs = append(cmdArgs, tt.args...)
			cmd.SetArgs(cmdArgs)

			// Capture output
			output := testhelpers.CaptureOutput(t, func() error {
				return cmd.Execute()
			})

			// Verify output contains expected text
			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Output does not contain expected text.\nGot: %s\nWant substring: %s", output, tt.wantOutput)
			}
		})
	}
}
