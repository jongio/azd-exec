package executor

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestDetectShell(t *testing.T) {
	tests := []struct {
		name       string
		scriptPath string
		want       string
	}{
		{
			name:       "PowerShell script",
			scriptPath: "test.ps1",
			want:       getPowerShellName(),
		},
		{
			name:       "Bash script",
			scriptPath: "test.sh",
			want:       "bash",
		},
		{
			name:       "Cmd script",
			scriptPath: "test.cmd",
			want:       "cmd",
		},
		{
			name:       "Batch script",
			scriptPath: "test.bat",
			want:       "cmd",
		},
		{
			name:       "Zsh script",
			scriptPath: "test.zsh",
			want:       "zsh",
		},
	}

	exec := New(Config{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exec.detectShell(tt.scriptPath)
			if got != tt.want {
				t.Errorf("detectShell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadShebang(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "Bash shebang",
			content: "#!/bin/bash\necho hello",
			want:    "bash",
		},
		{
			name:    "Sh shebang",
			content: "#!/bin/sh\necho hello",
			want:    "sh",
		},
		{
			name:    "Python shebang",
			content: "#!/usr/bin/env python3\nprint('hello')",
			want:    "python3",
		},
		{
			name:    "No shebang",
			content: "echo hello",
			want:    "",
		},
	}

	exec := New(Config{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			tmpFile, err := os.CreateTemp("", "script-*.sh")
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = os.Remove(tmpFile.Name())
			}()

			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatal(err)
			}
			_ = tmpFile.Close()

			got := exec.readShebang(tmpFile.Name())
			if got != tt.want {
				t.Errorf("readShebang() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name       string
		scriptName string
		content    string
		shell      string
		wantErr    bool
	}{
		{
			name:       "Simple bash script",
			scriptName: "test.sh",
			content:    "#!/bin/bash\necho 'Hello from script'\nexit 0",
			shell:      "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip on Windows if bash not available
			if runtime.GOOS == "windows" && tt.shell == "" {
				t.Skip("Skipping bash test on Windows")
			}

			// Create temp script
			tmpDir := t.TempDir()
			scriptPath := filepath.Join(tmpDir, tt.scriptName)
			// #nosec G306 - Script files need execute permission to run
			if err := os.WriteFile(scriptPath, []byte(tt.content), 0o700); err != nil {
				t.Fatal(err)
			}

			exec := New(Config{
				Shell: tt.shell,
			})

			err := exec.Execute(context.Background(), scriptPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func getPowerShellName() string {
	if runtime.GOOS == "windows" {
		return "powershell"
	}
	return "pwsh"
}
