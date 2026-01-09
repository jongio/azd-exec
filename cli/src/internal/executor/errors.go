package executor

import "fmt"

// ValidationError indicates that input validation failed.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Reason)
}

// ScriptNotFoundError indicates that a script file could not be found.
type ScriptNotFoundError struct {
	Path string
}

func (e *ScriptNotFoundError) Error() string {
	return fmt.Sprintf("script not found: %s", e.Path)
}

// InvalidShellError indicates that an invalid shell was specified.
type InvalidShellError struct {
	Shell string
}

func (e *InvalidShellError) Error() string {
	return fmt.Sprintf("invalid shell: %s (valid: bash, sh, zsh, pwsh, powershell, cmd)", e.Shell)
}

// ExecutionError indicates that script execution failed with an exit code.
type ExecutionError struct {
	ExitCode int
	Shell    string
	IsInline bool
}

func (e *ExecutionError) Error() string {
	if e.IsInline {
		return fmt.Sprintf("inline script exited with code %d (shell: %s)", e.ExitCode, e.Shell)
	}
	return fmt.Sprintf("script exited with code %d (shell: %s)", e.ExitCode, e.Shell)
}
