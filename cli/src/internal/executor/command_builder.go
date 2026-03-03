package executor

import (
	"os/exec"
	"strings"

	"github.com/jongio/azd-core/shellutil"
)

// validShells maps known shell names to whether they're valid.
// This is used for validation and shell-specific argument construction.
var validShells = map[string]bool{
	shellutil.ShellBash:       true,
	shellutil.ShellSh:         true,
	shellutil.ShellZsh:        true,
	shellutil.ShellPwsh:       true,
	shellutil.ShellPowerShell: true,
	shellutil.ShellCmd:        true,
}

// buildCommand builds the exec.Cmd for the given shell and script.
// It constructs shell-specific argument lists:
//   - Unix shells (bash, sh, zsh): Use -c for inline, direct path for files
//   - PowerShell (pwsh, powershell): Use -Command for inline, -File for files
//   - Windows cmd: Use /c for both inline and files
//   - Unknown shells: Fall back to -c flag (Unix-like behavior)
//
// Known shell names are normalized to lowercase for the executable binary
// to ensure correct lookup on case-sensitive filesystems.
// Script arguments (e.config.Args) are appended after the script specification.
func (e *Executor) buildCommand(shell, scriptOrPath string, isInline bool) *exec.Cmd {
	var cmdArgs []string
	skipAppendArgs := false

	// Normalize shell name to lowercase for both comparison and executable name.
	// This ensures correct binary lookup on case-sensitive filesystems
	// (e.g., --shell BASH resolves to "bash", not "BASH").
	shellLower := strings.ToLower(shell)

	// Use the lowercase name for known shells; keep original for unknown shells
	// (custom interpreters like "Python3" should preserve user's casing).
	shellBin := shell
	if validShells[shellLower] {
		shellBin = shellLower
	}

	switch shellLower {
	case shellutil.ShellBash, shellutil.ShellSh, shellutil.ShellZsh:
		if isInline {
			cmdArgs = []string{shellBin, "-c", scriptOrPath}
		} else {
			cmdArgs = []string{shellBin, scriptOrPath}
		}
	case shellutil.ShellPwsh, shellutil.ShellPowerShell:
		if isInline {
			cmdArgs = []string{shellBin, "-Command", e.buildPowerShellInlineCommand(scriptOrPath)}
			skipAppendArgs = true
		} else {
			cmdArgs = []string{shellBin, "-File", scriptOrPath}
		}
	case shellutil.ShellCmd:
		cmdArgs = []string{shellBin, "/c", scriptOrPath}
	default:
		// Unknown shell: use Unix-like -c pattern as fallback.
		// Preserve original casing for custom interpreters.
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
