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
	shell        string
	interactive  bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exec [script-file] [-- script-args...]",
		Short: "Execute scripts with Azure Developer CLI context",
		Long: `Execute scripts with full access to the Azure Developer CLI context.
Scripts will have access to all azd environment variables, configuration,
and can use azd commands within the script.

Examples:
	azd exec ./setup.sh
	azd exec ./deploy.ps1 --shell pwsh
	azd exec ./build.sh -- --verbose --config release
	azd exec ./init.sh -i  # Interactive mode
	
	azd exec version  # Show version
`,
		Args: func(cmd *cobra.Command, args []string) error {
			// Allow subcommands (version, listen) or script execution
			if len(args) == 0 {
				return fmt.Errorf("requires a script file or subcommand")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Execute script (default behavior)
			return executeScript(args)
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

	// Add script execution flags
	rootCmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run script in interactive mode")

	// Add extension-specific flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Output format: default or json")

	// Add azd global flags
	// These flags match the global flags available in azd to ensure compatibility
	// Without these, the extension will error when users pass global flags like --debug or --no-prompt
	rootCmd.PersistentFlags().StringVarP(&environment, "environment", "e", "", "The name of the environment to use.")
	rootCmd.PersistentFlags().StringVarP(&cwd, "cwd", "C", "", "Sets the current working directory.")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enables debugging and diagnostics logging.")
	rootCmd.PersistentFlags().BoolVar(&noPrompt, "no-prompt", false, "Accepts the default value instead of prompting, or it fails if there is no default.")
	rootCmd.PersistentFlags().StringVar(&traceLogFile, "trace-log-file", "", "Write a diagnostics trace to a file.")
	rootCmd.PersistentFlags().StringVar(&traceLogURL, "trace-log-url", "", "Send traces to an Open Telemetry compatible endpoint.")

	// Mark trace flags as hidden since they're advanced debugging features
	_ = rootCmd.PersistentFlags().MarkHidden("trace-log-file")
	_ = rootCmd.PersistentFlags().MarkHidden("trace-log-url")

	// Register subcommands
	rootCmd.AddCommand(
		commands.NewVersionCommand(&outputFormat),
		commands.NewListenCommand(), // Required for azd extension framework
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func executeScript(args []string) error {
	scriptPath := args[0]
	var scriptArgs []string
	if len(args) > 1 {
		scriptArgs = args[1:]
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
		WorkingDir:  cwd,
		Interactive: interactive,
		Args:        scriptArgs,
	})

	// Execute the script
	return exec.Execute(context.Background(), absPath)
}

// init remaps legacy invocation `azd script` to `azd exec` for compatibility.
func init() {
	if len(os.Args) > 1 && os.Args[1] == "script" {
		os.Args[1] = "exec"
	}
}
