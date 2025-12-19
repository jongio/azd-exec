//go:build integration
// +build integration

package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunCommandIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get path to test projects
	testProjectsDir := getTestProjectsDir(t)

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
			output := captureOutput(t, func() error {
				return cmd.Execute()
			})

			// Verify output contains expected text
			if !strings.Contains(output, tt.wantOutput) {
				t.Errorf("Output does not contain expected text.\nGot: %s\nWant substring: %s", output, tt.wantOutput)
			}
		})
	}
}

func TestRunCommandIntegration_InvalidScript(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)
	cmd.SetArgs([]string{"nonexistent-script.sh"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for nonexistent script, got nil")
	}
}

func TestRunCommandIntegration_WorkingDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if runtime.GOOS == "windows" {
		t.Skip("Skipping bash test on Windows")
	}

	testProjectsDir := getTestProjectsDir(t)
	scriptPath := filepath.Join(testProjectsDir, "bash", "simple.sh")

	// Verify script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Fatalf("Test script not found: %s", scriptPath)
	}

	// Create temp directory for working dir test
	tmpDir := t.TempDir()

	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)
	cmd.Flags().Set("cwd", tmpDir)
	cmd.SetArgs([]string{scriptPath})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() with working directory failed: %v", err)
	}
}

// Helper function to get test projects directory
func getTestProjectsDir(t *testing.T) string {
	t.Helper()

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Try multiple possible paths
	possiblePaths := []string{
		filepath.Join(cwd, "..", "..", "..", "tests", "projects"),       // From commands dir
		filepath.Join(cwd, "..", "tests", "projects"),                   // From src dir
		filepath.Join(cwd, "tests", "projects"),                         // From cli dir
		filepath.Join(cwd, "..", "..", "..", "..", "tests", "projects"), // From nested test dir
	}

	for _, testDir := range possiblePaths {
		testDir = filepath.Clean(testDir)
		if _, err := os.Stat(testDir); err == nil {
			return testDir
		}
	}

	// If not found, try to find the cli directory
	testDir := filepath.Join(cwd, "tests", "projects")
	for i := 0; i < 5; i++ {
		testDir = filepath.Join("..", testDir)
		testDir = filepath.Clean(testDir)
		if _, err := os.Stat(testDir); err == nil {
			return testDir
		}
	}

	t.Fatalf("Test projects directory not found. CWD: %s", cwd)
	return ""
}

// Helper function to capture command output
func captureOutput(t *testing.T, fn func() error) string {
	t.Helper()

	// Save original stdout
	origStdout := os.Stdout

	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Replace stdout
	os.Stdout = w

	// Channel for output
	outCh := make(chan string, 1)
	go func() {
		var output strings.Builder
		buf := make([]byte, 1024)
		for {
			n, readErr := r.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if readErr != nil {
				break
			}
		}
		outCh <- output.String()
	}()

	// Execute function
	fnErr := fn()

	// Close write end and restore stdout
	w.Close()
	os.Stdout = origStdout

	// Get output
	output := <-outCh

	if fnErr != nil {
		t.Logf("Command error: %v", fnErr)
	}

	return output
}
