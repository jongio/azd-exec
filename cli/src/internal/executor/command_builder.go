package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// validShells maps known shell names to whether they're valid.
var validShells = map[string]bool{
	shellBash:       true,
	shellSh:         true,
	shellZsh:        true,
	shellPwsh:       true,
	shellPowerShell: true,
	shellCmd:        true,
}

// buildCommand builds the exec.Cmd for the given shell and script.
func (e *Executor) buildCommand(shell, scriptOrPath string, isInline bool) *exec.Cmd {
	// Validate shell parameter
	if !validShells[strings.ToLower(shell)] {
		fmt.Fprintf(os.Stderr, "Warning: unknown shell '%s', execution may fail\n", shell)
	}

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
