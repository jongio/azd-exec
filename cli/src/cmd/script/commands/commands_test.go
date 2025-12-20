package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jongio/azd-exec/cli/src/internal/testhelpers"
)

func TestNewRunCommand(t *testing.T) {
	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)

	if cmd == nil {
		t.Fatal("NewRunCommand returned nil")
	}

	if !testhelpers.Contains(cmd.Use, "run") {
		t.Errorf("Command Use = %v, should contain 'run'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}

func TestRunCommandNoArgs(t *testing.T) {
	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error when no script provided, got nil")
	}
}

func TestRunCommandScriptNotFound(t *testing.T) {
	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)
	cmd.SetArgs([]string{"nonexistent-script.sh"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for nonexistent script, got nil")
	}
	if !testhelpers.Contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestRunCommandWithValidScript(t *testing.T) {
	// Create temporary script file
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.sh")
	content := []byte("#!/bin/bash\necho 'test'\n")
	if err := os.WriteFile(scriptPath, content, 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)
	cmd.SetArgs([]string{scriptPath})

	// Verify command accepts valid script path
	// Actual execution will depend on platform, so we allow common exec errors
	err := cmd.Execute()
	if err != nil {
		// Allow shell-not-found or path-related errors, but not arg parsing errors
		if !testhelpers.Contains(err.Error(), "executable file not found") &&
			!testhelpers.Contains(err.Error(), "exec format error") &&
			!testhelpers.Contains(err.Error(), "No such file or directory") &&
			!testhelpers.Contains(err.Error(), "code 127") {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

func TestRunCommandWithShellFlag(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.ps1")
	content := []byte("Write-Host 'test'\n")
	if err := os.WriteFile(scriptPath, content, 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	outputFormat := "default"
	cmd := NewRunCommand(&outputFormat)
	cmd.SetArgs([]string{"--shell", "pwsh", scriptPath})

	// Test that command accepts the flag without error in parsing
	err := cmd.Execute()
	// Execution may fail if pwsh not available, but flag should parse correctly
	if err != nil && testhelpers.Contains(err.Error(), "unknown flag") {
		t.Errorf("Shell flag not recognized: %v", err)
	}
}

func TestNewVersionCommand(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)

	if cmd == nil {
		t.Fatal("NewVersionCommand returned nil")
	}

	if cmd.Use != "version" {
		t.Errorf("Command Use = %v, want version", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}
}

func TestVersionCommandDefault(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)

	// Version command uses fmt.Printf which goes to stdout
	// We test by executing and checking error only
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully
}

func TestVersionCommandQuiet(t *testing.T) {
	outputFormat := "default"
	cmd := NewVersionCommand(&outputFormat)
	cmd.SetArgs([]string{"--quiet"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully with quiet flag
}

func TestVersionCommandJSON(t *testing.T) {
	outputFormat := "json"
	cmd := NewVersionCommand(&outputFormat)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}
	// Command executed successfully with JSON format
	// Note: Output validation would require capturing stdout which is complex in tests
}

func TestNewListenCommand(t *testing.T) {
	cmd := NewListenCommand()

	if cmd == nil {
		t.Fatal("NewListenCommand returned nil")
	}

	if cmd.Use != "listen" {
		t.Errorf("Command Use = %v, want listen", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}

	if !cmd.Hidden {
		t.Error("Listen command should be hidden")
	}
}

func TestListenCommandExecution(t *testing.T) {
	cmd := NewListenCommand()

	// Listen command should execute without error (it's a placeholder)
	err := cmd.Execute()
	if err != nil {
		t.Errorf("Listen command should not error, got: %v", err)
	}
}
