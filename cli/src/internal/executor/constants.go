package executor

// Shell identifiers.
const (
	shellBash       = "bash"
	shellCmd        = "cmd"
	shellPowerShell = "powershell"
	shellPwsh       = "pwsh"
	shellSh         = "sh"
	shellZsh        = "zsh"
)

// Operating system identifiers.
const (
	osWindows = "windows"
)

// File reading constants.
const (
	// shebangPrefix is the expected start of a shebang line.
	shebangPrefix = "#!"

	// shebangReadSize is the number of bytes to read for shebang detection.
	shebangReadSize = 2

	// envCommand is the common env wrapper in shebangs.
	envCommand = "env"
)
