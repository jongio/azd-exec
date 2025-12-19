package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jongio/azd-exec/cli/src/internal/executor"
	"github.com/spf13/cobra"
)

// NewRunCommand creates a new run command.
func NewRunCommand() *cobra.Command {
	var (
		shell       string
		workingDir  string
		interactive bool
		args        []string
	)

		cmd := &cobra.Command{
				Use:   "run [script-file] [-- script-args...]",
				Short: "Execute a script file with azd context",
				Long: `Execute a script file with access to the Azure Developer CLI context.
The script will have access to all azd environment variables, configuration,
and can use azd commands within the script.

Examples:
	azd exec run ./setup.sh
	azd exec run ./deploy.ps1 --shell pwsh
	azd exec run ./build.sh -- --verbose --config release
	azd exec run ./init.sh -i  # Interactive mode
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, cmdArgs []string) error {
			scriptPath := cmdArgs[0]

			// Handle script arguments after --
			if len(cmdArgs) > 1 {
				args = cmdArgs[1:]
			}

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
				Shell:       shell,
				WorkingDir:  workingDir,
				Interactive: interactive,
				Args:        args,
			})

			// Execute the script
			return exec.Execute(context.Background(), absPath)
		},
	}

	cmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.")
	cmd.Flags().StringVarP(&workingDir, "cwd", "C", "", "Working directory for script execution (defaults to script directory)")
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run script in interactive mode")

	return cmd
}
