package executor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteWithWorkingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.bat")

	// Create a simple batch script
	content := "@echo off\necho test\nexit 0"
	if err := os.WriteFile(scriptPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	// Create executor with working directory
	exec := New(Config{
		WorkingDir: tmpDir,
	})

	// Verify config
	if exec.config.WorkingDir != tmpDir {
		t.Errorf("WorkingDir = %v, want %v", exec.config.WorkingDir, tmpDir)
	}
}

func TestExecuteWithInteractiveMode(t *testing.T) {
	exec := New(Config{
		Interactive: true,
	})

	if !exec.config.Interactive {
		t.Error("Interactive mode not set")
	}
}

func TestExecuteWithCustomShell(t *testing.T) {
	exec := New(Config{
		Shell: "bash",
	})

	if exec.config.Shell != "bash" {
		t.Errorf("Shell = %v, want bash", exec.config.Shell)
	}
}

func TestExecuteError(t *testing.T) {
	// Try to execute non-existent script
	exec := New(Config{})

	err := exec.Execute(context.Background(), "nonexistent-script.sh")
	if err == nil {
		t.Error("Expected error for nonexistent script, got nil")
	}
}

func TestExecuteWithArguments(t *testing.T) {
	args := []string{"arg1", "arg2", "arg3"}
	exec := New(Config{Args: args})

	if len(exec.config.Args) != 3 {
		t.Errorf("Args length = %v, want 3", len(exec.config.Args))
	}

	if exec.config.Args[0] != "arg1" {
		t.Errorf("Args[0] = %v, want arg1", exec.config.Args[0])
	}
}
