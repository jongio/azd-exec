package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jongio/azd-exec/cli/src/cmd/script/commands"
	"github.com/jongio/azd-exec/cli/src/internal/executor"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	debugMode    bool
	noPrompt     bool
	cwd          string
	environment  string
	traceLogFile string
	traceLogURL  string
	// Root command flags for direct script execution
	shell       string
	workingDir  string
	interactive bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exec [script-file] [-- script-args...]",
		Short: "Exec - Execute commands/scripts with Azure Developer CLI context",
		Long: `Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.

Examples:
	azd exec ./setup.sh
	azd exec ./deploy.ps1 --shell pwsh
	azd exec ./build.sh -- --verbose --config release
	azd exec ./init.sh -i  # Interactive mode`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			// Resolve script path
			scriptPath := args[0]
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
				Args:        []string{}, // TODO: Handle args after --
			})

			// Execute the script
			return exec.Execute(context.Background(), absPath)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Handle working directory change
			if cwd != "" {
				if err := os.Chdir(cwd); err != nil {
					return fmt.Errorf("failed to change working directory to %s: %w", cwd, err)
				}
			}

			// Handle debug mode
			if debugMode {
				_ = os.Setenv("AZD_SCRIPT_DEBUG", "true")
			}

			// Handle no-prompt mode
			if noPrompt {
				_ = os.Setenv("AZD_NO_PROMPT", "true")
			}

			// Handle environment selection
			if environment != "" {
				_ = os.Setenv("AZURE_ENV_NAME", environment)
			}

			// Handle trace logging
			if traceLogFile != "" {
				_ = os.Setenv("AZD_TRACE_LOG_FILE", traceLogFile)
			}
			if traceLogURL != "" {
				_ = os.Setenv("AZD_TRACE_LOG_URL", traceLogURL)
			}

			return nil
		},
	}

	// Add extension-specific flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Output format: default or json")

	// Add flags for direct script execution (when using 'azd exec ./script.sh')
	rootCmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.")
	rootCmd.Flags().StringVarP(&workingDir, "working-dir", "w", "", "Working directory for script execution (defaults to script directory)")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run script in interactive mode")

	// Add azd global flags
	// These flags match the global flags available in azd to ensure compatibility
	// Without these, the extension will error when users pass global flags like --debug or --no-prompt
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&noPrompt, "no-prompt", false, "Disable prompts")
	rootCmd.PersistentFlags().StringVarP(&cwd, "cwd", "C", "", "Sets the current working directory")
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "The name of the environment to use")
	rootCmd.PersistentFlags().StringVar(&traceLogFile, "trace-log-file", "", "Write a diagnostics trace to a file.")
	rootCmd.PersistentFlags().StringVar(&traceLogURL, "trace-log-url", "", "Send traces to an Open Telemetry compatible endpoint.")

	// Mark trace flags as hidden since they're advanced debugging features
	_ = rootCmd.PersistentFlags().MarkHidden("trace-log-file")
	_ = rootCmd.PersistentFlags().MarkHidden("trace-log-url")

	// Register subcommands
	rootCmd.AddCommand(
		commands.NewVersionCommand(&outputFormat),
		commands.NewListenCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init remaps legacy invocation `azd script` to `azd exec` for compatibility.
func init() {
	if len(os.Args) > 1 && os.Args[1] == "script" {
		os.Args[1] = "exec"
	}
}
