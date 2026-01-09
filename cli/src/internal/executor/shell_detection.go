package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// detectShell auto-detects the appropriate shell based on the script extension and shebang.
// Detection priority:
//  1. File extension (.ps1, .cmd, .bat, .sh, .zsh)
//  2. Shebang line (#!/bin/bash, #!/usr/bin/env python3, etc.)
//  3. OS-specific default (Windows: cmd, Unix: bash)
//
// Returns the shell command name (e.g., "bash", "pwsh", "cmd").
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
		// Check shebang line for scripts without recognized extensions
		if shebang := e.readShebang(scriptPath); shebang != "" {
			return shebang
		}

		// Default based on OS (cmd for Windows, bash for Unix)
		if runtime.GOOS == osWindows {
			return shellCmd
		}
		return shellBash
	}
}

// readShebang reads the shebang line from a script file and extracts the shell name.
// It handles common shebang formats:
//   - #!/bin/bash
//   - #!/usr/bin/env python3
//   - #! /bin/sh
//
// Returns:
//   - Empty string if no shebang is found or file cannot be read
//   - The base name of the shell/interpreter (e.g., "bash", "python3")
func (e *Executor) readShebang(scriptPath string) string {
	file, err := os.Open(scriptPath) // #nosec G304 - scriptPath is validated by caller
	if err != nil {
		return ""
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			// Log error but don't fail - we may have already read what we needed
			// Only log to stderr if we're in debug mode to avoid noise
			if os.Getenv(envVarScriptDebug) == "true" {
				fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", filepath.Base(scriptPath), closeErr)
			}
		}
	}()

	reader := bufio.NewReader(file)

	// Read first bytes to check for shebang
	buf := make([]byte, shebangReadSize)
	if _, readErr := io.ReadFull(reader, buf); readErr != nil {
		return ""
	}

	if string(buf) != shebangPrefix {
		return ""
	}

	// Read the rest of the line
	line, lineErr := reader.ReadString('\n')
	if lineErr != nil && lineErr != io.EOF {
		return ""
	}

	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return ""
	}

	// Handle "#!/usr/bin/env python3" style shebangs
	if filepath.Base(parts[0]) == envCommand && len(parts) > 1 {
		return filepath.Base(parts[1])
	}

	return filepath.Base(parts[0])
}
