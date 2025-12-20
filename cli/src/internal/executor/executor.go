package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
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

	// Build command
	cmd := e.buildCommand(shell, scriptPath)
	cmd.Dir = workingDir

	// Inherit all environment variables (includes azd context)
	envVars := os.Environ()

	// Resolve Key Vault references if any exist
	if e.hasKeyVaultReferences(envVars) {
		resolvedEnvVars, err := e.resolveKeyVaultReferences(ctx, envVars)
		if err != nil {
			// Log warning but continue with unresolved variables
			// This allows scripts to run even if Key Vault resolution fails
			fmt.Fprintf(os.Stderr, "Warning: Failed to resolve Key Vault references: %v\n", err)
			fmt.Fprintf(os.Stderr, "Continuing with original environment variables...\n")
		} else {
			envVars = resolvedEnvVars
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
	if os.Getenv("AZD_SCRIPT_DEBUG") == "true" {
		fmt.Fprintf(os.Stderr, "Executing: %s %s\n", shell, strings.Join(cmd.Args[1:], " "))
		fmt.Fprintf(os.Stderr, "Working directory: %s\n", workingDir)
	}

	// Run the command
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("script exited with code %d", exitErr.ExitCode())
		}
		return fmt.Errorf("failed to execute script %q with shell %q: %w", filepath.Base(scriptPath), shell, err)
	}

	return nil
}

// detectShell auto-detects the appropriate shell based on the script extension.
func (e *Executor) detectShell(scriptPath string) string {
	ext := strings.ToLower(filepath.Ext(scriptPath))

	switch ext {
	case ".ps1":
		if runtime.GOOS == osWindows {
			return shellPowerShell
		}
		return shellPwsh
	case ".cmd", ".bat":
		return shellCmd
	case ".sh":
		return shellBash
	case ".zsh":
		return shellZsh
	default:
		// Check shebang line
		if shebang := e.readShebang(scriptPath); shebang != "" {
			return shebang
		}

		// Default based on OS
		if runtime.GOOS == osWindows {
			return shellCmd
		}
		return shellBash
	}
}

// readShebang reads the shebang line from a script file.
func (e *Executor) readShebang(scriptPath string) string {
	file, err := os.Open(scriptPath) // #nosec G304 - scriptPath is validated by caller
	if err != nil {
		return ""
	}
	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)

	// Read first bytes to check for shebang
	buf := make([]byte, shebangReadSize)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return ""
	}

	if string(buf) != shebangPrefix {
		return ""
	}

	// Read the rest of the line
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return ""
	}

	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) > 0 {
		// Handle "#!/usr/bin/env python3" style shebangs
		if filepath.Base(parts[0]) == envCommand && len(parts) > 1 {
			return filepath.Base(parts[1])
		}
		shellPath := parts[0]
		return filepath.Base(shellPath)
	}

	return ""
}

// buildCommand builds the exec.Cmd for the given shell and script.
func (e *Executor) buildCommand(shell, scriptPath string) *exec.Cmd {
	var cmdArgs []string

	switch strings.ToLower(shell) {
	case shellBash, shellSh, shellZsh:
		cmdArgs = []string{shell, scriptPath}
	case shellPwsh, shellPowerShell:
		cmdArgs = []string{shell, "-File", scriptPath}
	case shellCmd:
		cmdArgs = []string{shell, "/c", scriptPath}
	default:
		cmdArgs = []string{shell, scriptPath}
	}

	// Append script arguments
	if len(e.config.Args) > 0 {
		cmdArgs = append(cmdArgs, e.config.Args...)
	}

	return exec.Command(cmdArgs[0], cmdArgs[1:]...) // #nosec G204 - cmdArgs are controlled by caller
}

// hasKeyVaultReferences checks if any environment variables contain Key Vault references.
func (e *Executor) hasKeyVaultReferences(envVars []string) bool {
	for _, envVar := range envVars {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 && IsKeyVaultReference(parts[1]) {
			return true
		}
	}
	return false
}

// resolveKeyVaultReferences resolves all Key Vault references in environment variables.
func (e *Executor) resolveKeyVaultReferences(ctx context.Context, envVars []string) ([]string, error) {
	resolver, err := NewKeyVaultResolver()
	if err != nil {
		return nil, fmt.Errorf("failed to create Key Vault resolver: %w", err)
	}

	return resolver.ResolveEnvironmentVariables(ctx, envVars)
}
