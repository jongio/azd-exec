package executor

import (
	"os/exec"
	"strings"
)

// buildCommand builds the exec.Cmd for the given shell and script.
func (e *Executor) buildCommand(shell, scriptOrPath string, isInline bool) *exec.Cmd {
	var cmdArgs []string

	switch strings.ToLower(shell) {
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
		if isInline {
			cmdArgs = []string{shell, "/c", scriptOrPath}
		} else {
			cmdArgs = []string{shell, "/c", scriptOrPath}
		}
	default:
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
