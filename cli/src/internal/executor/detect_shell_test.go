package executor

import (
	"os"
	"runtime"
	"testing"
)

func TestDetectShellUnit(t *testing.T) {
	// Get expected PowerShell name based on OS
	expectedPwsh := "powershell"
	if runtime.GOOS != "windows" {
		expectedPwsh = "pwsh"
	}

	tests := []struct {
		name       string
		scriptPath string
		want       string
	}{
		{
			name:       "PowerShell script",
			scriptPath: "test.ps1",
			want:       expectedPwsh,
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
		{
			name:       "No extension",
			scriptPath: "script",
			want: func() string {
				if runtime.GOOS == "windows" {
					return "cmd"
				}
				return "bash"
			}(),
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

func TestReadShebangUnit(t *testing.T) {
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
			name:    "Zsh shebang",
			content: "#!/usr/bin/zsh\necho hello",
			want:    "zsh",
		},
		{
			name:    "No shebang",
			content: "echo hello",
			want:    "",
		},
		{
			name:    "Empty file",
			content: "",
			want:    "",
		},
		{
			name:    "Shebang with spaces",
			content: "#! /bin/bash\necho hello",
			want:    "bash",
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

func TestGetPowerShellNameUnit(t *testing.T) {
	// Test PowerShell name selection based on OS
	exec := New(Config{})
	got := exec.detectShell("test.ps1")

	if runtime.GOOS == "windows" {
		if got != "powershell" {
			t.Errorf("detectShell(.ps1) on Windows = %v, want powershell", got)
		}
	} else {
		if got != "pwsh" {
			t.Errorf("detectShell(.ps1) on non-Windows = %v, want pwsh", got)
		}
	}
}
