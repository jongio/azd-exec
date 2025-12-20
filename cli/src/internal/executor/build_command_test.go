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
		wantShell  string
		wantArgs   int // minimum number of args expected
	}{
		{
			name:       "Bash with no args",
			shell:      "bash",
			scriptPath: "test.sh",
			args:       nil,
			wantShell:  "bash",
			wantArgs:   1, // script path
		},
		{
			name:       "Bash with args",
			shell:      "bash",
			scriptPath: "test.sh",
			args:       []string{"arg1", "arg2"},
			wantShell:  "bash",
			wantArgs:   3, // script path + 2 args
		},
		{
			name:       "PowerShell with args",
			shell:      expectedPwsh,
			scriptPath: "test.ps1",
			args:       []string{"-Param1", "value"},
			wantShell:  expectedPwsh,
			wantArgs:   4, // shell + -File + script path + args
		},
		{
			name:       "Cmd with args",
			shell:      "cmd",
			scriptPath: "test.bat",
			args:       []string{"arg1"},
			wantShell:  "cmd",
			wantArgs:   3, // shell + /c + script path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(Config{Shell: tt.shell, Args: tt.args})
			cmd := exec.buildCommand(tt.shell, tt.scriptPath)

			if cmd.Path == "" {
				t.Error("buildCommand() returned command with empty Path")
			}

			if len(cmd.Args) < tt.wantArgs {
				t.Errorf("buildCommand() args count = %v, want at least %v\nArgs: %v", len(cmd.Args), tt.wantArgs, cmd.Args)
			}
		})
	}
}
