//go:build integration
// +build integration

package commands

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jongio/azd-exec/cli/src/internal/testhelpers"
)

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

	testProjectsDir := testhelpers.GetTestProjectsDir(t)
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
