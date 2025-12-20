package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jongio/azd-exec/cli/src/internal/executor"
)

// ScriptOptions contains options for script execution.
type ScriptOptions struct {
	Shell       string
	WorkingDir  string
	Interactive bool
	Args        []string
}

// ExecuteScript executes a script file with the given options.
func ExecuteScript(ctx context.Context, scriptPath string, opts ScriptOptions) error {
	// Resolve script path
	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("failed to resolve script path: %w", err)
	}

	// Check if script exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("script file not found: %s", scriptPath)
	}

	// Create executor
	exec := executor.New(executor.Config{
		Shell:       opts.Shell,
		WorkingDir:  opts.WorkingDir,
		Interactive: opts.Interactive,
		Args:        opts.Args,
	})

	// Execute the script
	return exec.Execute(ctx, absPath)
}
