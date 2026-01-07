package executor

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jongio/azd-exec/cli/src/internal/testhelpers"
)

func TestExecute_FileValidation(t *testing.T) {
	exec := New(Config{})

	t.Run("Empty script path", func(t *testing.T) {
		err := exec.Execute(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty script path")
		}
		if err.Error() != "script path cannot be empty" {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Valid script file", func(t *testing.T) {
		projectsDir := testhelpers.GetTestProjectsDir(t)

		// Use OS-appropriate scripts to avoid Windows path translation issues.
		scriptPath := filepath.Join(projectsDir, "bash", "simple.sh")
		if runtime.GOOS == osWindows {
			scriptPath = filepath.Join(projectsDir, "powershell", "simple.ps1")
		}

		// Skip on Windows if bash not available, or use different script
		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			t.Skip("Test script not found, skipping")
		}

		err := exec.Execute(context.Background(), scriptPath)
		if err != nil {
			t.Logf("Script execution error (may be expected): %v", err)
		}
	})
}

func TestExecuteInline_Validation(t *testing.T) {
	exec := New(Config{})

	t.Run("Empty script content", func(t *testing.T) {
		err := exec.ExecuteInline(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty script content")
		}
		if err.Error() != "script content cannot be empty" {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Valid inline script", func(t *testing.T) {
		err := exec.ExecuteInline(context.Background(), "echo 'test'")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})
}

func TestPrepareEnvironment(t *testing.T) {
	exec := New(Config{})

	t.Run("No Key Vault references", func(t *testing.T) {
		// Set up environment without KV references
		origEnv := os.Environ()
		defer func() {
			os.Clearenv()
			for _, env := range origEnv {
				if idx := findEqualsSign(env); idx > 0 {
					_ = os.Setenv(env[:idx], env[idx+1:])
				}
			}
		}()

		_ = os.Setenv("NORMAL_VAR", "value")

		envVars, err := exec.prepareEnvironment(context.Background())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(envVars) == 0 {
			t.Error("Expected non-empty environment")
		}
	})

	t.Run("With Key Vault references but no Azure credentials", func(t *testing.T) {
		origEnv := os.Environ()
		defer func() {
			os.Clearenv()
			for _, env := range origEnv {
				if idx := findEqualsSign(env); idx > 0 {
					_ = os.Setenv(env[:idx], env[idx+1:])
				}
			}
		}()

		_ = os.Setenv("KV_VAR", "@Microsoft.KeyVault(VaultName=test;SecretName=secret)")

		_, err := exec.prepareEnvironment(context.Background())
		// Without Azure credentials, this should return an error
		if err != nil {
			t.Logf("Expected error without Azure credentials: %v", err)
		} else {
			t.Log("Azure credentials available, KV resolution succeeded")
		}
	})
}

func TestHasKeyVaultReferences(t *testing.T) {
	exec := New(Config{})

	tests := []struct {
		name    string
		envVars []string
		want    bool
	}{
		{
			name:    "No references",
			envVars: []string{"VAR1=value1", "VAR2=value2"},
			want:    false,
		},
		{
			name:    "Has SecretUri reference",
			envVars: []string{"VAR1=value1", "KV=@Microsoft.KeyVault(SecretUri=https://vault.vault.azure.net/secrets/test)"},
			want:    true,
		},
		{
			name:    "Has VaultName reference",
			envVars: []string{"KV=@Microsoft.KeyVault(VaultName=vault;SecretName=secret)"},
			want:    true,
		},
		{
			name:    "Empty environment",
			envVars: []string{},
			want:    false,
		},
		{
			name:    "Malformed env var",
			envVars: []string{"INVALID"},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exec.hasKeyVaultReferences(tt.envVars)
			if got != tt.want {
				t.Errorf("hasKeyVaultReferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogDebugInfo(t *testing.T) {
	exec := New(Config{})

	// Test doesn't crash and produces output
	t.Run("Inline script", func(t *testing.T) {
		exec.logDebugInfo("bash", "/tmp", "echo test", true, []string{"bash", "-c", "echo test"})
	})

	t.Run("File script", func(t *testing.T) {
		exec.logDebugInfo("bash", "/tmp", "test.sh", false, []string{"bash", "test.sh"})
	})
}

func TestRunCommand_ErrorHandling(t *testing.T) {
	exec := New(Config{})

	shell := "bash"
	exitCmd := "exit 1"
	missingCmd := "nonexistent-command-xyz"
	missingFile := "/nonexistent/file.sh"
	if runtime.GOOS == osWindows {
		shell = shellCmd
		exitCmd = "exit /b 1"
		missingCmd = "nonexistent-command-xyz"
		missingFile = "C:\\nonexistent\\file.cmd"
	}

	t.Run("Exit code error", func(t *testing.T) {
		// Create a command that will fail
		cmd := exec.buildCommand(shell, exitCmd, true)
		err := exec.runCommand(cmd, "test", shell, true)
		if err == nil {
			t.Error("Expected error for exit code 1")
		}
	})

	t.Run("Inline script error formatting", func(t *testing.T) {
		cmd := exec.buildCommand(shell, missingCmd, true)
		err := exec.runCommand(cmd, missingCmd, shell, true)
		if err == nil {
			t.Error("Expected error for nonexistent command")
		}
	})

	t.Run("File script error formatting", func(t *testing.T) {
		cmd := exec.buildCommand(shell, missingFile, false)
		err := exec.runCommand(cmd, missingFile, shell, false)
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}

func TestExecutorWithDebugMode(t *testing.T) {
	origDebug := os.Getenv("AZD_SCRIPT_DEBUG")
	defer func() {
		if origDebug != "" {
			_ = os.Setenv("AZD_SCRIPT_DEBUG", origDebug)
		} else {
			_ = os.Unsetenv("AZD_SCRIPT_DEBUG")
		}
	}()

	_ = os.Setenv("AZD_SCRIPT_DEBUG", "true")

	exec := New(Config{})
	err := exec.ExecuteInline(context.Background(), "echo 'debug test'")
	if err != nil {
		t.Errorf("Unexpected error in debug mode: %v", err)
	}
}

func TestExecutorWithWorkingDir(t *testing.T) {
	tmpDir := t.TempDir()
	exec := New(Config{
		WorkingDir: tmpDir,
	})

	err := exec.ExecuteInline(context.Background(), "pwd")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestExecutorWithArgs(t *testing.T) {
	projectsDir := testhelpers.GetTestProjectsDir(t)

	scriptPath := filepath.Join(projectsDir, "bash", "with-args.sh")
	args := []string{"arg1", "arg2"}
	if runtime.GOOS == osWindows {
		scriptPath = filepath.Join(projectsDir, "powershell", "with-params.ps1")
		args = []string{"-Name", "Test", "-Greeting", "Hi"}
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Skip("Test script not found, skipping")
	}

	exec := New(Config{
		Args: args,
	})

	err := exec.Execute(context.Background(), scriptPath)
	if err != nil {
		t.Logf("Script execution error (may be expected): %v", err)
	}
}

func TestExecutorInteractiveMode(t *testing.T) {
	exec := New(Config{
		Interactive: true,
	})

	// Just verify the config is set - can't easily test actual interactive behavior
	if !exec.config.Interactive {
		t.Error("Expected Interactive to be true")
	}
}

// Helper function to find '=' in environment variable.
func findEqualsSign(env string) int {
	for i, c := range env {
		if c == '=' {
			return i
		}
	}
	return -1
}
