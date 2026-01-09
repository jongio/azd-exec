package executor

// Shell identifiers used for script execution.
// These constants define the supported shell types and are used
// for shell detection and command building.
const (
	// shellBash is the Bourne Again Shell (default on most Unix systems).
	shellBash = "bash"

	// shellCmd is the Windows Command Prompt.
	shellCmd = "cmd"

	// shellPowerShell is Windows PowerShell (5.1 and earlier).
	shellPowerShell = "powershell"

	// shellPwsh is PowerShell Core (6.0+, cross-platform).
	shellPwsh = "pwsh"

	// shellSh is the POSIX shell.
	shellSh = "sh"

	// shellZsh is the Z Shell.
	shellZsh = "zsh"
)

// Operating system identifiers.
const (
	// osWindows identifies the Windows operating system.
	osWindows = "windows"
)

// Environment variable names.
const (
	// envVarScriptDebug enables debug output for script execution.
	// When set to "true", execution details are logged to stderr.
	envVarScriptDebug = "AZD_SCRIPT_DEBUG"
)

// File reading constants for shebang detection.
const (
	// shebangPrefix is the expected start of a shebang line ("#!").
	shebangPrefix = "#!"

	// shebangReadSize is the number of bytes to read for shebang detection.
	// This must be at least len(shebangPrefix) bytes.
	shebangReadSize = len(shebangPrefix)

	// envCommand is the common env wrapper in shebangs (e.g., #!/usr/bin/env bash).
	envCommand = "env"
)
