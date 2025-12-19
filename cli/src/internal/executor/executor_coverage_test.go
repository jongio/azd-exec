package executor

import (
	"context"
	"os"
	"os/exec"
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

func TestDetectShellWithShebang(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name     string
		content  string
		filename string
		want     string
	}{
		{
			name:     "Bash shebang overrides .txt extension",
			content:  "#!/bin/bash\necho hello",
			filename: "script.txt",
			want:     "bash",
		},
		{
			name:     "Python shebang",
			content:  "#!/usr/bin/env python3\nprint('hello')",
			filename: "script",
			want:     "python3",
		},
		{
			name:     "Zsh shebang",
			content:  "#!/usr/bin/zsh\necho hello",
			filename: "script",
			want:     "zsh",
		},
	}
	
	exec := New(Config{})
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(scriptPath, []byte(tt.content), 0600); err != nil {
				t.Fatal(err)
			}
			
			got := exec.detectShell(scriptPath)
			if got != tt.want {
				t.Errorf("detectShell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildCommandWithCustomShell(t *testing.T) {
	tests := []struct {
		name       string
		shell      string
		scriptPath string
		args       []string
		wantFirst  string
	}{
		{
			name:       "Custom shell python",
			shell:      "python3",
			scriptPath: "script.py",
			args:       []string{"arg1"},
			wantFirst:  "python3",
		},
		{
			name:       "Custom shell node",
			shell:      "node",
			scriptPath: "script.js",
			args:       nil,
			wantFirst:  "node",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(Config{Shell: tt.shell, Args: tt.args})
			cmd := exec.buildCommand(tt.shell, tt.scriptPath)
			
			// Check if command was built
			if cmd == nil {
				t.Fatal("buildCommand returned nil")
			}
			
			// Verify args contain the shell
			found := false
			for _, arg := range cmd.Args {
				if arg == tt.wantFirst {
					found = true
					break
				}
			}
			
			if !found {
				t.Errorf("buildCommand args don't contain shell %v: %v", tt.wantFirst, cmd.Args)
			}
		})
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

func TestBuildCommandShellVariations(t *testing.T) {
	tests := []struct {
		shell      string
		scriptPath string
		wantArgs   []string
	}{
		{
			shell:      "SH",
			scriptPath: "test.sh",
			wantArgs:   []string{"sh", "test.sh"},
		},
		{
			shell:      "BASH",
			scriptPath: "test.sh",
			wantArgs:   []string{"bash", "test.sh"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			exec := New(Config{Shell: tt.shell})
			cmd := exec.buildCommand(tt.shell, tt.scriptPath)
			
			if cmd == nil {
				t.Fatal("buildCommand returned nil")
			}
		})
	}
}

func TestReadShebangFileNotFound(t *testing.T) {
	exec := New(Config{})
	got := exec.readShebang("nonexistent-file.sh")
	
	if got != "" {
		t.Errorf("readShebang(nonexistent) = %v, want empty string", got)
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

func TestBuildCommandLookPath(t *testing.T) {
	// Test that buildCommand creates a valid exec.Cmd
	exec := New(Config{})
	cmd := exec.buildCommand("cmd", "test.bat")
	
	// On Windows, cmd should be findable
	if cmd.Path == "" {
		t.Error("buildCommand created command with empty Path")
	}
	
	// Verify we can look up the command
	_, err := execLookPath(cmd.Args[0])
	if err != nil {
		t.Logf("Command %v not found in PATH (may be platform-specific)", cmd.Args[0])
	}
}

// Helper to wrap exec.LookPath for testing
func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}
