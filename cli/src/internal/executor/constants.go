package executor

import "github.com/jongio/azd-core/shellutil"

// Shell identifiers used for script execution.
// These constants define the supported shell types and are used
// for shell detection and command building.
const (
	// shellBash is the Bourne Again Shell (default on most Unix systems).
	shellBash = shellutil.ShellBash

	// shellCmd is the Windows Command Prompt.
	shellCmd = shellutil.ShellCmd

	// shellPowerShell is Windows PowerShell (5.1 and earlier).
	shellPowerShell = shellutil.ShellPowerShell

	// shellPwsh is PowerShell Core (6.0+, cross-platform).
	shellPwsh = shellutil.ShellPwsh

	// shellSh is the POSIX shell.
	shellSh = shellutil.ShellSh

	// shellZsh is the Z Shell.
	shellZsh = shellutil.ShellZsh
)

// Operating system identifiers.
const (
	// osWindows identifies the Windows operating system.
	osWindows = "windows"
)
