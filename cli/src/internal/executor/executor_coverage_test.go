package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jongio/azd-core/testutil"
)

func TestExecute_FileValidation(t *testing.T) {
	exec := New(Config{})

	t.Run("Empty script path", func(t *testing.T) {
		err := exec.Execute(context.Background(), "")
		if err == nil {
			t.Error("Expected error for empty script path")
		}
		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Valid script file", func(t *testing.T) {
		projectsDir := testutil.FindTestData(t, "tests", "projects")

		// Use OS-appropriate scripts to avoid Windows path translation issues.
		scriptPath := filepath.Join(projectsDir, "bash", "simple.sh")
		if runtime.GOOS == "windows" {
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
		if !strings.Contains(err.Error(), "cannot be empty") {
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
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					_ = os.Setenv(parts[0], parts[1])
				}
			}
		}()

		_ = os.Setenv("NORMAL_VAR", "value")

		envVars, _, err := exec.prepareEnvironment(context.Background())
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// warnings are allowed but should be empty for non-KV env
		if len(envVars) == 0 {
			t.Error("Expected non-empty environment")
		}
	})

	t.Run("With Key Vault references but no Azure credentials", func(t *testing.T) {
		origEnv := os.Environ()
		defer func() {
			os.Clearenv()
			for _, env := range origEnv {
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 {
					_ = os.Setenv(parts[0], parts[1])
				}
			}
		}()

		_ = os.Setenv("KV_VAR", "@Microsoft.KeyVault(VaultName=test;SecretName=secret)")

		_, warnings, err := exec.prepareEnvironment(context.Background())
		// Default behavior is continue-on-error; errors should be reserved for strict mode.
		if err != nil {
			t.Fatalf("Unexpected error in continue-on-error mode: %v", err)
		}
		if len(warnings) == 0 {
			t.Log("No warnings emitted (credentials+secret may have resolved successfully)")
		} else {
			t.Logf("Warnings emitted for Key Vault resolution as expected: %d", len(warnings))
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
	if runtime.GOOS == "windows" {
		shell = "cmd"
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

func TestExecutorWithArgs(t *testing.T) {
	projectsDir := testutil.FindTestData(t, "tests", "projects")

	scriptPath := filepath.Join(projectsDir, "bash", "with-args.sh")
	args := []string{"arg1", "arg2"}
	if runtime.GOOS == "windows" {
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

// TestExecute_DirectoryPath verifies that executing a directory returns an error.
func TestExecute_DirectoryPath(t *testing.T) {
	exec := New(Config{})
	err := exec.Execute(context.Background(), t.TempDir())
	if err == nil {
		t.Error("Expected error for directory path")
	}
	if !strings.Contains(err.Error(), "must be a file") {
		t.Errorf("Unexpected error: %v", err)
	}
}

// TestPrepareEnvironment_StopOnKeyVaultError verifies fail-fast behavior.
func TestPrepareEnvironment_StopOnKeyVaultError(t *testing.T) {
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range origEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	}()

	os.Clearenv()
	_ = os.Setenv("KV_VAR", "@Microsoft.KeyVault(VaultName=test;SecretName=secret)")

	exec := New(Config{StopOnKeyVaultError: true})
	_, _, err := exec.prepareEnvironment(context.Background())
	// With StopOnKeyVaultError=true and no real Azure credentials, we expect an error
	// (either resolver creation fails or resolution fails)
	if err != nil {
		t.Logf("Got expected error in fail-fast mode: %v", err)
	} else {
		t.Log("Resolver succeeded (Azure credentials available)")
	}
}

// TestPrepareEnvironment_ResolverCreationError verifies handling when Key Vault resolver fails to initialize.
func TestPrepareEnvironment_ResolverCreationError(t *testing.T) {
	origEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, env := range origEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	}()

	os.Clearenv()
	_ = os.Setenv("KV_VAR", "@Microsoft.KeyVault(VaultName=test;SecretName=secret)")

	// Inject a resolver that always fails
	oldResolver := newKeyVaultEnvResolver
	defer func() { newKeyVaultEnvResolver = oldResolver }()
	newKeyVaultEnvResolver = func() (keyVaultEnvResolver, error) {
		return nil, fmt.Errorf("mock resolver error")
	}

	t.Run("continue on error mode", func(t *testing.T) {
		exec := New(Config{StopOnKeyVaultError: false})
		envVars, warnings, err := exec.prepareEnvironment(context.Background())
		if err != nil {
			t.Fatalf("expected no error in continue mode, got: %v", err)
		}
		if len(warnings) == 0 {
			t.Error("expected warnings when resolver fails")
		}
		if len(envVars) == 0 {
			t.Error("expected fallback to original env vars")
		}
	})

	t.Run("stop on error mode", func(t *testing.T) {
		exec := New(Config{StopOnKeyVaultError: true})
		_, _, err := exec.prepareEnvironment(context.Background())
		if err == nil {
			t.Error("expected error in stop-on-error mode")
		}
		if !strings.Contains(err.Error(), "mock resolver error") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

// TestGetDefaultShellForOS verifies platform-specific shell detection.
func TestGetDefaultShellForOS(t *testing.T) {
	shell := getDefaultShellForOS()
	if shell == "" {
		t.Error("getDefaultShellForOS returned empty string")
	}
	if runtime.GOOS == "windows" {
		if shell != "powershell" {
			t.Errorf("expected powershell on Windows, got %q", shell)
		}
	} else {
		if shell != "bash" {
			t.Errorf("expected bash on non-Windows, got %q", shell)
		}
	}
}

// TestNewExecutor tests the New() constructor.
// Merged from constructor_test.go.
func TestNewExecutor(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name:   "Default config",
			config: Config{},
		},
		{
			name: "With shell specified",
			config: Config{
				Shell: "bash",
			},
		},
		{
			name: "With interactive mode",
			config: Config{
				Interactive: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := New(tt.config)
			if exec == nil {
				t.Error("New() returned nil")
				return
			}
			if exec.config.Shell != tt.config.Shell {
				t.Errorf("Shell = %v, want %v", exec.config.Shell, tt.config.Shell)
			}
			if exec.config.Interactive != tt.config.Interactive {
				t.Errorf("Interactive = %v, want %v", exec.config.Interactive, tt.config.Interactive)
			}
		})
	}
}
