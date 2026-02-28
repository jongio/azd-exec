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

func TestQuotePowerShellArg(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{name: "empty string", arg: "", want: "''"},
		{name: "simple arg", arg: "hello", want: "'hello'"},
		{name: "arg with single quote", arg: "it's", want: "'it''s'"},
		{name: "arg with multiple quotes", arg: "a'b'c", want: "'a''b''c'"},
		{name: "arg with double dash", arg: "--skip-sync", want: "'--skip-sync'"},
		{name: "arg with spaces", arg: "hello world", want: "'hello world'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := quotePowerShellArg(tt.arg)
			if got != tt.want {
				t.Errorf("quotePowerShellArg(%q) = %q, want %q", tt.arg, got, tt.want)
			}
		})
	}
}

func TestBuildPowerShellInlineCommand(t *testing.T) {
	t.Run("no args returns script as-is", func(t *testing.T) {
		e := New(Config{})
		got := e.buildPowerShellInlineCommand("Get-Date")
		if got != "Get-Date" {
			t.Errorf("got %q, want %q", got, "Get-Date")
		}
	})

	t.Run("with args joins and quotes", func(t *testing.T) {
		e := New(Config{Args: []string{"arg1", "it's"}})
		got := e.buildPowerShellInlineCommand("cmd")
		want := "cmd 'arg1' 'it''s'"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
