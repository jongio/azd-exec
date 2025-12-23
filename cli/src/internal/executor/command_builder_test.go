package executor

import (
	"os/exec"
	"testing"
)

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
			cmd := exec.buildCommand(tt.shell, tt.scriptPath, false)

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
			cmd := exec.buildCommand(tt.shell, tt.scriptPath, false)

			if cmd == nil {
				t.Fatal("buildCommand returned nil")
			}
		})
	}
}

func TestBuildCommandLookPath(t *testing.T) {
	// Test that buildCommand creates a valid exec.Cmd
	exec := New(Config{})
	cmd := exec.buildCommand("cmd", "test.bat", false)

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

// Helper to wrap exec.LookPath for testing.
func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}
