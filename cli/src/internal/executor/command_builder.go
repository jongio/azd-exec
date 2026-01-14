package executor

import (
	"os/exec"
	"strings"
)

// validShells maps known shell names to whether they're valid.
// This is used for validation and shell-specific argument construction.
var validShells = map[string]bool{
	shellBash:       true,
	shellSh:         true,
	shellZsh:        true,
	shellPwsh:       true,
	shellPowerShell: true,
	shellCmd:        true,
}

// buildCommand builds the exec.Cmd for the given shell and script.
// It constructs shell-specific argument lists:
//   - Unix shells (bash, sh, zsh): Use -c for inline, direct path for files
//   - PowerShell (pwsh, powershell): Use -Command for inline, -File for files
//   - Windows cmd: Use /c for both inline and files
//   - Unknown shells: Fall back to -c flag (Unix-like behavior)
//
// Script arguments (e.config.Args) are appended after the script specification.
func (e *Executor) buildCommand(shell, scriptOrPath string, isInline bool) *exec.Cmd {
	var cmdArgs []string
	skipAppendArgs := false

	// Normalize shell name to lowercase for comparison
	shellLower := strings.ToLower(shell)

	switch shellLower {
	case shellBash, shellSh, shellZsh:
		if isInline {
			cmdArgs = []string{shell, "-c", scriptOrPath}
		} else {
			cmdArgs = []string{shell, scriptOrPath}
		}
	case shellPwsh, shellPowerShell:
		if isInline {
			cmdArgs = []string{shell, "-Command", e.buildPowerShellInlineCommand(scriptOrPath)}
			skipAppendArgs = true
		} else {
			cmdArgs = []string{shell, "-File", scriptOrPath}
		}
	case shellCmd:
		cmdArgs = []string{shell, "/c", scriptOrPath}
	default:
		// Unknown shell: use Unix-like -c pattern as fallback
		if isInline {
			cmdArgs = []string{shell, "-c", scriptOrPath}
		} else {
			cmdArgs = []string{shell, scriptOrPath}
		}
	}

	// Append script arguments unless already embedded
	if !skipAppendArgs && len(e.config.Args) > 0 {
		cmdArgs = append(cmdArgs, e.config.Args...)
	}

	return exec.Command(cmdArgs[0], cmdArgs[1:]...) // #nosec G204 - cmdArgs are controlled by caller
}

// buildPowerShellInlineCommand joins the inline script with its arguments into a single
// -Command string to avoid PowerShell re-quoting passthrough arguments (e.g., "--flag").
// All arguments are single-quoted with internal quotes escaped to preserve literal values.
func (e *Executor) buildPowerShellInlineCommand(scriptOrPath string) string {
	if len(e.config.Args) == 0 {
		return scriptOrPath
	}

	quotedArgs := make([]string, len(e.config.Args))
	for i, arg := range e.config.Args {
		quotedArgs[i] = quotePowerShellArg(arg)
	}

	return strings.Join(append([]string{scriptOrPath}, quotedArgs...), " ")
}

// quotePowerShellArg returns a safely single-quoted PowerShell argument.
// Single quotes inside the argument are escaped by doubling them.
func quotePowerShellArg(arg string) string {
	if arg == "" {
		return "''"
	}

	return "'" + strings.ReplaceAll(arg, "'", "''") + "'"
}
