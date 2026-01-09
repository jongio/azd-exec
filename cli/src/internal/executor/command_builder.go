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
			cmdArgs = []string{shell, "-Command", scriptOrPath}
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

	// Append script arguments
	if len(e.config.Args) > 0 {
		cmdArgs = append(cmdArgs, e.config.Args...)
	}

	return exec.Command(cmdArgs[0], cmdArgs[1:]...) // #nosec G204 - cmdArgs are controlled by caller
}
