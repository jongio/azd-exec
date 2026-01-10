package executor

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jongio/azd-core/shellutil"
)

func TestBuildCommand_Bash(t *testing.T) {
	exec := New(Config{})

	t.Run("Inline bash", func(t *testing.T) {
		cmd := exec.buildCommand("bash", "echo test", true)
		if len(cmd.Args) < 3 {
			t.Fatalf("Expected at least 3 args, got %d", len(cmd.Args))
		}
		if cmd.Args[0] != "bash" {
			t.Errorf("Expected 'bash', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "-c" {
			t.Errorf("Expected '-c', got %q", cmd.Args[1])
		}
		if cmd.Args[2] != "echo test" {
			t.Errorf("Expected 'echo test', got %q", cmd.Args[2])
		}
	})

	t.Run("File bash", func(t *testing.T) {
		cmd := exec.buildCommand("bash", "script.sh", false)
		if len(cmd.Args) < 2 {
			t.Fatalf("Expected at least 2 args, got %d", len(cmd.Args))
		}
		if cmd.Args[0] != "bash" {
			t.Errorf("Expected 'bash', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "script.sh" {
			t.Errorf("Expected 'script.sh', got %q", cmd.Args[1])
		}
	})
}

func TestBuildCommand_PowerShell(t *testing.T) {
	exec := New(Config{})

	t.Run("Inline pwsh", func(t *testing.T) {
		cmd := exec.buildCommand("pwsh", "Write-Host 'test'", true)
		if len(cmd.Args) < 3 {
			t.Fatalf("Expected at least 3 args, got %d", len(cmd.Args))
		}
		if cmd.Args[0] != "pwsh" {
			t.Errorf("Expected 'pwsh', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "-Command" {
			t.Errorf("Expected '-Command', got %q", cmd.Args[1])
		}
	})

	t.Run("File powershell", func(t *testing.T) {
		cmd := exec.buildCommand("powershell", "script.ps1", false)
		if len(cmd.Args) < 3 {
			t.Fatalf("Expected at least 3 args, got %d", len(cmd.Args))
		}
		if cmd.Args[0] != "powershell" {
			t.Errorf("Expected 'powershell', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "-File" {
			t.Errorf("Expected '-File', got %q", cmd.Args[1])
		}
		if cmd.Args[2] != "script.ps1" {
			t.Errorf("Expected 'script.ps1', got %q", cmd.Args[2])
		}
	})
}

func TestBuildCommand_Cmd(t *testing.T) {
	exec := New(Config{})

	t.Run("Inline cmd", func(t *testing.T) {
		cmd := exec.buildCommand("cmd", "echo test", true)
		if len(cmd.Args) < 3 {
			t.Fatalf("Expected at least 3 args, got %d", len(cmd.Args))
		}
		if cmd.Args[0] != "cmd" {
			t.Errorf("Expected 'cmd', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "/c" {
			t.Errorf("Expected '/c', got %q", cmd.Args[1])
		}
	})

	t.Run("File cmd", func(t *testing.T) {
		cmd := exec.buildCommand("cmd", "script.bat", false)
		if len(cmd.Args) < 3 {
			t.Fatalf("Expected at least 3 args, got %d", len(cmd.Args))
		}
		if cmd.Args[1] != "/c" {
			t.Errorf("Expected '/c', got %q", cmd.Args[1])
		}
	})
}

func TestBuildCommand_WithArgs(t *testing.T) {
	exec := New(Config{
		Args: []string{"--verbose", "--output=json"},
	})

	cmd := exec.buildCommand("bash", "script.sh", false)
	if len(cmd.Args) < 4 {
		t.Fatalf("Expected at least 4 args, got %d", len(cmd.Args))
	}

	// Last two should be the script args
	if cmd.Args[len(cmd.Args)-2] != "--verbose" {
		t.Errorf("Expected '--verbose', got %q", cmd.Args[len(cmd.Args)-2])
	}
	if cmd.Args[len(cmd.Args)-1] != "--output=json" {
		t.Errorf("Expected '--output=json', got %q", cmd.Args[len(cmd.Args)-1])
	}
}

func TestBuildCommand_Zsh(t *testing.T) {
	exec := New(Config{})

	t.Run("Inline zsh", func(t *testing.T) {
		cmd := exec.buildCommand("zsh", "echo test", true)
		if cmd.Args[0] != "zsh" {
			t.Errorf("Expected 'zsh', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "-c" {
			t.Errorf("Expected '-c', got %q", cmd.Args[1])
		}
	})

	t.Run("File zsh", func(t *testing.T) {
		cmd := exec.buildCommand("zsh", "script.zsh", false)
		if cmd.Args[0] != "zsh" {
			t.Errorf("Expected 'zsh', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "script.zsh" {
			t.Errorf("Expected 'script.zsh', got %q", cmd.Args[1])
		}
	})
}

func TestBuildCommand_Sh(t *testing.T) {
	exec := New(Config{})

	cmd := exec.buildCommand("sh", "script.sh", false)
	if cmd.Args[0] != "sh" {
		t.Errorf("Expected 'sh', got %q", cmd.Args[0])
	}
}

func TestBuildCommand_CustomShell(t *testing.T) {
	exec := New(Config{})

	t.Run("Inline custom shell", func(t *testing.T) {
		cmd := exec.buildCommand("python3", "print('test')", true)
		if cmd.Args[0] != "python3" {
			t.Errorf("Expected 'python3', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "-c" {
			t.Errorf("Expected '-c', got %q", cmd.Args[1])
		}
	})

	t.Run("File custom shell", func(t *testing.T) {
		cmd := exec.buildCommand("python3", "script.py", false)
		if cmd.Args[0] != "python3" {
			t.Errorf("Expected 'python3', got %q", cmd.Args[0])
		}
		if cmd.Args[1] != "script.py" {
			t.Errorf("Expected 'script.py', got %q", cmd.Args[1])
		}
	})
}

func TestBuildCommand_CaseInsensitiveShellNames(t *testing.T) {
	exec := New(Config{})

	tests := []struct {
		shell    string
		expected string
	}{
		{"BASH", "BASH"},
		{"Bash", "Bash"},
		{"PWSH", "PWSH"},
		{"PowerShell", "PowerShell"},
		{"CMD", "CMD"},
	}

	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			cmd := exec.buildCommand(tt.shell, "test", true)
			if cmd.Args[0] != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, cmd.Args[0])
			}
		})
	}
}

func TestReadShebang_ValidShebang(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "bash shebang",
			content:  "#!/bin/bash\necho 'test'",
			expected: "bash",
		},
		{
			name:     "python shebang",
			content:  "#!/usr/bin/python3\nprint('test')",
			expected: "python3",
		},
		{
			name:     "env shebang",
			content:  "#!/usr/bin/env python3\nprint('test')",
			expected: "python3",
		},
		{
			name:     "sh shebang",
			content:  "#!/bin/sh\necho 'test'",
			expected: "sh",
		},
		{
			name:     "zsh shebang",
			content:  "#!/bin/zsh\necho 'test'",
			expected: "zsh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := filepath.Join(tmpDir, tt.name)
			if err := os.WriteFile(scriptPath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			result := shellutil.ReadShebang(scriptPath)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReadShebang_NoShebang(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "no shebang",
			content: "echo 'test'",
		},
		{
			name:    "comment not shebang",
			content: "# This is a comment\necho 'test'",
		},
		{
			name:    "empty file",
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := filepath.Join(tmpDir, tt.name)
			if err := os.WriteFile(scriptPath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			result := shellutil.ReadShebang(scriptPath)
			if result != "" {
				t.Errorf("Expected empty string, got %q", result)
			}
		})
	}
}

func TestReadShebang_FileErrors(t *testing.T) {
	t.Run("Nonexistent file", func(t *testing.T) {
		result := shellutil.ReadShebang("/nonexistent/file.sh")
		if result != "" {
			t.Errorf("Expected empty string for nonexistent file, got %q", result)
		}
	})
}

// TestBuildCommand tests buildCommand with various shells and scenarios.
// Merged from build_command_test.go.
func TestBuildCommand(t *testing.T) {
	// Get expected PowerShell name based on OS
	expectedPwsh := "powershell"
	if runtime.GOOS != "windows" {
		expectedPwsh = "pwsh"
	}

	tests := []struct {
		name       string
		shell      string
		scriptPath string
		args       []string
		isInline   bool
		wantShell  string
		wantArgs   int // minimum number of args expected
	}{
		{
			name:       "Bash with no args (file)",
			shell:      "bash",
			scriptPath: "test.sh",
			args:       nil,
			isInline:   false,
			wantShell:  "bash",
			wantArgs:   1, // script path
		},
		{
			name:       "Bash with args (file)",
			shell:      "bash",
			scriptPath: "test.sh",
			args:       []string{"arg1", "arg2"},
			isInline:   false,
			wantShell:  "bash",
			wantArgs:   3, // script path + 2 args
		},
		{
			name:       "Bash inline script",
			shell:      "bash",
			scriptPath: "echo hello",
			args:       nil,
			isInline:   true,
			wantShell:  "bash",
			wantArgs:   2, // -c + script content
		},
		{
			name:       "PowerShell with args (file)",
			shell:      expectedPwsh,
			scriptPath: "test.ps1",
			args:       []string{"-Param1", "value"},
			isInline:   false,
			wantShell:  expectedPwsh,
			wantArgs:   4, // shell + -File + script path + args
		},
		{
			name:       "PowerShell inline script",
			shell:      expectedPwsh,
			scriptPath: "Write-Host 'hello'",
			args:       nil,
			isInline:   true,
			wantShell:  expectedPwsh,
			wantArgs:   2, // -Command + script content
		},
		{
			name:       "Cmd with args (file)",
			shell:      "cmd",
			scriptPath: "test.bat",
			args:       []string{"arg1"},
			isInline:   false,
			wantShell:  "cmd",
			wantArgs:   3, // shell + /c + script path
		},
		{
			name:       "Cmd inline script",
			shell:      "cmd",
			scriptPath: "echo hello",
			args:       nil,
			isInline:   true,
			wantShell:  "cmd",
			wantArgs:   2, // /c + script content
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(Config{Shell: tt.shell, Args: tt.args})
			cmd := exec.buildCommand(tt.shell, tt.scriptPath, tt.isInline)

			if cmd.Path == "" {
				t.Error("buildCommand() returned command with empty Path")
			}

			if len(cmd.Args) < tt.wantArgs {
				t.Errorf("buildCommand() args count = %v, want at least %v\nArgs: %v", len(cmd.Args), tt.wantArgs, cmd.Args)
			}
		})
	}
}
