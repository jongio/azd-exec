package executor

import (
	"runtime"
	"testing"
)

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
