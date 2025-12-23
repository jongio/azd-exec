package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Config holds the configuration for script execution.
type Config struct {
	Shell       string   // Shell to use for execution
	WorkingDir  string   // Working directory
	Interactive bool     // Interactive mode
	Args        []string // Arguments to pass to the script
}

// Executor executes scripts with azd context.
type Executor struct {
	config Config
}

// New creates a new script executor.
func New(config Config) *Executor {
	return &Executor{config: config}
}

// Execute runs a script file with azd context.
func (e *Executor) Execute(ctx context.Context, scriptPath string) error {
	// Validate script path
	if scriptPath == "" {
		return fmt.Errorf("script path cannot be empty")
	}

	// Auto-detect shell if not specified
	shell := e.config.Shell
	if shell == "" {
		shell = e.detectShell(scriptPath)
	}

	// Determine working directory
	workingDir := e.config.WorkingDir
	if workingDir == "" {
		workingDir = filepath.Dir(scriptPath)
	}

	return e.executeCommand(ctx, shell, workingDir, scriptPath, false)
}

// ExecuteInline runs an inline script command with azd context.
func (e *Executor) ExecuteInline(ctx context.Context, scriptContent string) error {
	// Validate script content
	if scriptContent == "" {
		return fmt.Errorf("script content cannot be empty")
	}

	// Auto-detect shell if not specified, default based on OS
	shell := e.config.Shell
	if shell == "" {
		if runtime.GOOS == osWindows {
			shell = shellPowerShell
		} else {
			shell = shellBash
		}
	}

	// Determine working directory
	workingDir := e.config.WorkingDir
	if workingDir == "" {
		var err error
		workingDir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	return e.executeCommand(ctx, shell, workingDir, scriptContent, true)
}

// executeCommand is the common execution logic for both file and inline scripts.
func (e *Executor) executeCommand(ctx context.Context, shell, workingDir, scriptOrPath string, isInline bool) error {
	// Build command
	cmd := e.buildCommand(shell, scriptOrPath, isInline)
	cmd.Dir = workingDir

	// Prepare environment with Key Vault resolution
	envVars, err := e.prepareEnvironment(ctx)
	if err != nil {
		// Log warning but continue with unresolved variables
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		fmt.Fprintf(os.Stderr, "Continuing with original environment variables...\n")
		envVars = os.Environ()
	}
	cmd.Env = envVars

	// Set up stdio
	if e.config.Interactive {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Add debug output
	if os.Getenv("AZD_SCRIPT_DEBUG") == "true" {
		e.logDebugInfo(shell, workingDir, scriptOrPath, isInline, cmd.Args)
	}

	// Run the command
	return e.runCommand(cmd, scriptOrPath, shell, isInline)
}

// prepareEnvironment prepares environment variables with Key Vault resolution.
func (e *Executor) prepareEnvironment(ctx context.Context) ([]string, error) {
	envVars := os.Environ()

	if !e.hasKeyVaultReferences(envVars) {
		return envVars, nil
	}

	resolver, err := NewKeyVaultResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Vault resolver: %w", err)
	}

	resolvedVars, err := resolver.ResolveEnvironmentVariables(ctx, envVars)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve Key Vault references: %w", err)
	}

	return resolvedVars, nil
}

// logDebugInfo logs debug information about script execution.
func (e *Executor) logDebugInfo(shell, workingDir, scriptOrPath string, isInline bool, cmdArgs []string) {
	if isInline {
		fmt.Fprintf(os.Stderr, "Executing inline: %s\n", shell)
		fmt.Fprintf(os.Stderr, "Script content: %s\n", scriptOrPath)
	} else {
		fmt.Fprintf(os.Stderr, "Executing: %s %s\n", shell, strings.Join(cmdArgs[1:], " "))
	}
	fmt.Fprintf(os.Stderr, "Working directory: %s\n", workingDir)
}

// runCommand executes the command and handles errors.
func (e *Executor) runCommand(cmd *exec.Cmd, scriptOrPath, shell string, isInline bool) error {
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("script exited with code %d", exitErr.ExitCode())
		}
		if isInline {
			return fmt.Errorf("failed to execute inline script with shell %q: %w", shell, err)
		}
		return fmt.Errorf("failed to execute script %q with shell %q: %w", filepath.Base(scriptOrPath), shell, err)
	}
	return nil
}

// hasKeyVaultReferences checks if any environment variables contain Key Vault references.
func (e *Executor) hasKeyVaultReferences(envVars []string) bool {
	for _, envVar := range envVars {
		if parts := strings.SplitN(envVar, "=", 2); len(parts) == 2 && IsKeyVaultReference(parts[1]) {
			return true
		}
	}
	return false
}
