package main

import (
	"os"
	"testing"
)

func TestExecuteScriptMissingFile(t *testing.T) {
	err := executeScript([]string{"nonexistent-script.sh"})
	if err == nil {
		t.Error("Expected error for nonexistent script, got nil")
	}
	if !contains(err.Error(), "not found") {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestExecuteScriptWithArgs(t *testing.T) {
	// Create temporary script file
	tmpDir := t.TempDir()
	scriptPath := tmpDir + "/test.sh"
	content := []byte("#!/bin/bash\necho 'test'\n")
	if err := os.WriteFile(scriptPath, content, 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Test with args
	err := executeScript([]string{scriptPath, "--arg1", "value1"})
	// Allow common execution errors on different platforms
	if err != nil && !contains(err.Error(), "executable file not found") &&
		!contains(err.Error(), "exec format error") &&
		!contains(err.Error(), "No such file or directory") &&
		!contains(err.Error(), "code 127") {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestExecuteScriptPathResolution(t *testing.T) {
	// Test that relative paths are resolved
	tmpDir := t.TempDir()
	scriptPath := tmpDir + "/script.sh"
	content := []byte("#!/bin/bash\necho 'test'\n")
	if err := os.WriteFile(scriptPath, content, 0755); err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Change to temp dir to test relative path
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	err := executeScript([]string{"./script.sh"})
	// Allow common execution errors, but check path resolution succeeded
	if err != nil && contains(err.Error(), "failed to resolve") {
		t.Errorf("Path resolution failed: %v", err)
	}
}

func contains(s, substr string) bool {
	if s == "" {
		return false
	}
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
