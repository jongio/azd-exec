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

// detectShell auto-detects the appropriate shell based on the script extension.
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
		// Check shebang line
		if shebang := e.readShebang(scriptPath); shebang != "" {
			return shebang
		}

		// Default based on OS
		if runtime.GOOS == osWindows {
			return shellCmd
		}
		return shellBash
	}
}

// readShebang reads the shebang line from a script file.
func (e *Executor) readShebang(scriptPath string) string {
	file, err := os.Open(scriptPath) // #nosec G304 - scriptPath is validated by caller
	if err != nil {
		return ""
	}
	defer func() {
		if err := file.Close(); err != nil {
			// Log error but don't fail - we may have already read what we needed
			fmt.Fprintf(os.Stderr, "warning: failed to close file %s: %v\n", scriptPath, err)
		}
	}()

	reader := bufio.NewReader(file)

	// Read first bytes to check for shebang
	buf := make([]byte, shebangReadSize)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return ""
	}

	if string(buf) != shebangPrefix {
		return ""
	}

	// Read the rest of the line
	line, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return ""
	}

	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) > 0 {
		// Handle "#!/usr/bin/env python3" style shebangs
		if filepath.Base(parts[0]) == envCommand && len(parts) > 1 {
			return filepath.Base(parts[1])
		}
		shellPath := parts[0]
		return filepath.Base(shellPath)
	}

	return ""
}
