//go:build integration
// +build integration

package executor

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestExecuteIntegration tests real script execution with various shells and scenarios.
func TestExecuteIntegration(t *testing.T) {
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

// TestExecuteInlineIntegration tests inline script execution with various shells.
func TestExecuteInlineIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name    string
		script  string
		shell   string
		wantErr bool
		skip    func() bool
	}{
		{
			name:    "Bash inline echo",
			script:  "echo 'Hello from inline'",
			shell:   "bash",
			wantErr: false,
			skip:    func() bool { return runtime.GOOS == "windows" },
		},
		{
			name:    "Bash inline with env var",
			script:  "echo $HOME",
			shell:   "bash",
			wantErr: false,
			skip:    func() bool { return runtime.GOOS == "windows" },
		},
		{
			name:    "PowerShell inline echo",
			script:  "Write-Host 'Hello from PowerShell'",
			shell:   getPowerShellName(),
			wantErr: false,
			skip:    func() bool { return false },
		},
		{
			name:    "PowerShell inline with env var",
			script:  "Write-Host $env:USERPROFILE",
			shell:   getPowerShellName(),
			wantErr: false,
			skip:    func() bool { return runtime.GOOS != "windows" },
		},
		{
			name:    "Bash multi-line inline",
			script:  "echo 'line 1'; echo 'line 2'",
			shell:   "bash",
			wantErr: false,
			skip:    func() bool { return runtime.GOOS == "windows" },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip() {
				t.Skip("Skipping test on this platform")
			}

			exec := New(Config{
				Shell: tt.shell,
			})

			err := exec.ExecuteInline(context.Background(), tt.script)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteInline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExecuteInlineEmptyScript tests that empty inline scripts return an error.
func TestExecuteInlineEmptyScript(t *testing.T) {
	exec := New(Config{})
	err := exec.ExecuteInline(context.Background(), "")
	if err == nil {
		t.Error("ExecuteInline() with empty script should return an error")
	}
}
