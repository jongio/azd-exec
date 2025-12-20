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
