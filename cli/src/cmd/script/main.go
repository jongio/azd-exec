package main

import (
	"fmt"
	"os"

	"github.com/jongio/azd-exec/cli/src/cmd/script/commands"
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
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exec",
		Short: "Exec - Execute commands/scripts with Azure Developer CLI context",
		Long:  `Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.`,
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

	// Register all commands
	rootCmd.AddCommand(
		commands.NewRunCommand(&outputFormat),
		commands.NewVersionCommand(&outputFormat),
		commands.NewListenCommand(), // Required for azd extension framework
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
