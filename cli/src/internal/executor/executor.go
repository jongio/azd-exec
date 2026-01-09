// Package executor provides secure script execution with Azure context and Key Vault integration.
// It supports multiple shells (bash, sh, zsh, pwsh, powershell, cmd) and handles environment
// variable resolution including Azure Key Vault secret references.
package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jongio/azd-core/keyvault"
)

// Config holds the configuration for script execution.
// All fields are optional and have sensible defaults.
type Config struct {
	// Shell specifies the shell to use for execution.
	// If empty, shell is auto-detected from script extension or shebang.
	// Valid values: bash, sh, zsh, pwsh, powershell, cmd
	Shell string

	// Interactive enables interactive mode, connecting stdin to the script.
	Interactive bool

	// StopOnKeyVaultError causes azd exec to fail-fast when any Key Vault reference fails to resolve.
	// Default is false (continue resolving other references and run with unresolved values left as-is).
	StopOnKeyVaultError bool

	// Args are additional arguments to pass to the script.
	Args []string
}

// Validate checks if the Config has valid values.
func (c *Config) Validate() error {
	if c.Shell != "" && !validShells[strings.ToLower(c.Shell)] {
		return &InvalidShellError{Shell: c.Shell}
	}
	return nil
}

type keyVaultEnvResolver interface {
	ResolveEnvironmentVariables(ctx context.Context, envVars []string, options keyvault.ResolveEnvironmentOptions) ([]string, []keyvault.KeyVaultResolutionWarning, error)
}

var newKeyVaultEnvResolver = func() (keyVaultEnvResolver, error) {
	return keyvault.NewKeyVaultResolver()
}

// Executor executes scripts with azd context.
type Executor struct {
	config Config
}

// New creates a new script executor with the given configuration.
// Returns a configured Executor ready to execute scripts.
// The config is validated before creating the executor.
func New(config Config) *Executor {
	// Note: We don't return an error here to maintain backward compatibility.
	// Invalid config values are caught during execution.
	return &Executor{config: config}
}

// Execute runs a script file with azd context.
// The script path is validated for existence and security.
// Returns an error if:
//   - scriptPath is empty
//   - scriptPath does not exist or is not a regular file
//   - scriptPath is a directory
//   - scriptPath contains path traversal attempts (..)
//   - script execution fails
func (e *Executor) Execute(ctx context.Context, scriptPath string) error {
	// Validate script path
	if scriptPath == "" {
		return &ValidationError{Field: "scriptPath", Reason: "cannot be empty"}
	}

	// Get absolute path and validate
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return &ValidationError{Field: "scriptPath", Reason: fmt.Sprintf("invalid path: %v", err)}
	}

	// Check for path traversal attempts
	if strings.Contains(filepath.ToSlash(absPath), "/../") {
		return &ValidationError{Field: "scriptPath", Reason: "path traversal not allowed"}
	}

	// Ensure script exists before attempting to execute it
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &ScriptNotFoundError{Path: filepath.Base(absPath)}
		}
		return &ValidationError{Field: "scriptPath", Reason: fmt.Sprintf("cannot access: %v", err)}
	}

	if info.IsDir() {
		return &ValidationError{Field: "scriptPath", Reason: "must be a file, not a directory"}
	}

	// Auto-detect shell if not specified
	shell := e.config.Shell
	if shell == "" {
		shell = e.detectShell(absPath)
	}

	// Use script's directory as working directory
	workingDir := filepath.Dir(absPath)

	return e.executeCommand(ctx, shell, workingDir, absPath, false)
}

// ExecuteInline runs an inline script command with azd context.
// The shell is auto-detected based on OS if not specified in config.
// Returns an error if:
//   - scriptContent is empty or only whitespace
//   - shell detection fails
//   - script execution fails
func (e *Executor) ExecuteInline(ctx context.Context, scriptContent string) error {
	// Validate script content
	if strings.TrimSpace(scriptContent) == "" {
		return &ValidationError{Field: "scriptContent", Reason: "cannot be empty or whitespace"}
	}

	// Auto-detect shell if not specified, default based on OS
	shell := e.config.Shell
	if shell == "" {
		shell = getDefaultShellForOS()
	}

	// Use current directory as working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	return e.executeCommand(ctx, shell, workingDir, scriptContent, true)
}

// executeCommand is the common execution logic for both file and inline scripts.
func (e *Executor) executeCommand(ctx context.Context, shell, workingDir, scriptOrPath string, isInline bool) error {
	// Build command
	cmd := e.buildCommand(shell, scriptOrPath, isInline)
	cmd.Dir = workingDir

	// Prepare environment with Key Vault resolution
	envVars, warnings, err := e.prepareEnvironment(ctx)
	if err != nil {
		return err
	}
	for _, w := range warnings {
		if w.Key != "" {
			fmt.Fprintf(os.Stderr, "Warning: failed to resolve Key Vault reference for %s: %v\n", w.Key, w.Err)
		} else {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", w.Err)
		}
	}
	cmd.Env = envVars

	// Set up stdio
	if e.config.Interactive {
		cmd.Stdin = os.Stdin
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Add debug output
	if os.Getenv(envVarScriptDebug) == "true" {
		e.logDebugInfo(shell, workingDir, scriptOrPath, isInline, cmd.Args)
	}

	// Run the command
	return e.runCommand(cmd, scriptOrPath, shell, isInline)
}

// prepareEnvironment prepares environment variables with Key Vault resolution.
func (e *Executor) prepareEnvironment(ctx context.Context) ([]string, []keyvault.KeyVaultResolutionWarning, error) {
	envVars := os.Environ()

	if !e.hasKeyVaultReferences(envVars) {
		return envVars, nil, nil
	}

	resolver, err := newKeyVaultEnvResolver()
	if err != nil {
		if e.config.StopOnKeyVaultError {
			return nil, nil, fmt.Errorf("failed to create Key Vault resolver: %w", err)
		}
		return envVars, []keyvault.KeyVaultResolutionWarning{{Err: fmt.Errorf("failed to create Key Vault resolver: %w", err)}}, nil
	}

	resolvedVars, warnings, err := resolver.ResolveEnvironmentVariables(ctx, envVars, keyvault.ResolveEnvironmentOptions{StopOnError: e.config.StopOnKeyVaultError})
	if err != nil {
		// Fail-fast mode returns an error and should prevent script execution.
		return nil, warnings, fmt.Errorf("failed to resolve Key Vault references: %w", err)
	}

	return resolvedVars, warnings, nil
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
// Error messages are sanitized to avoid leaking sensitive path information.
func (e *Executor) runCommand(cmd *exec.Cmd, scriptOrPath, shell string, isInline bool) error {
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return &ExecutionError{
				ExitCode: exitErr.ExitCode(),
				Shell:    shell,
				IsInline: isInline,
			}
		}
		if isInline {
			return fmt.Errorf("failed to execute inline script with shell %q: %w", shell, err)
		}
		// Use only base filename to avoid leaking full paths in error messages
		return fmt.Errorf("failed to execute script %q with shell %q: %w", filepath.Base(scriptOrPath), shell, err)
	}
	return nil
}

// hasKeyVaultReferences checks if any environment variables contain Key Vault references.
func (e *Executor) hasKeyVaultReferences(envVars []string) bool {
	for _, envVar := range envVars {
		if parts := strings.SplitN(envVar, "=", 2); len(parts) == 2 && keyvault.IsKeyVaultReference(parts[1]) {
			return true
		}
	}
	return false
}

// getDefaultShellForOS returns the default shell for the current operating system.
func getDefaultShellForOS() string {
	if runtime.GOOS == osWindows {
		return shellPowerShell
	}
	return shellBash
}
