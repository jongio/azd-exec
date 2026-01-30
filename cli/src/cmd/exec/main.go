// Package main provides the azd exec extension command-line interface.
// It enables execution of scripts with full Azure Developer CLI context,
// including environment variables and Azure credentials.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jongio/azd-core/cliout"
	"github.com/jongio/azd-core/env"
	"github.com/jongio/azd-exec/cli/src/cmd/exec/commands"
	"github.com/jongio/azd-exec/cli/src/internal/executor"
	"github.com/spf13/cobra"
)

var (
	// Output and logging flags.
	outputFormat string
	debugMode    bool
	noPrompt     bool

	// Execution context flags.
	cwd         string
	environment string

	// Tracing flags (advanced debugging).
	traceLogFile string
	traceLogURL  string

	// Root command flags for direct script execution.
	shell       string
	interactive bool

	// Key Vault resolution behavior flags.
	stopOnKeyVaultError bool
)

type scriptExecutor interface {
	Execute(ctx context.Context, scriptPath string) error
	ExecuteInline(ctx context.Context, scriptContent string) error
}

var newScriptExecutor = func(config executor.Config) scriptExecutor {
	return executor.New(config)
}

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		cliout.Error("%v", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "exec [script-file-or-command] [-- script-args...]",
		Short: "Exec - Execute commands/scripts with Azure Developer CLI context",
		Long: `Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.

Examples:
\tazd exec ./setup.sh                           # Execute script file
\tazd exec 'echo $AZURE_ENV_NAME'               # Inline bash (Linux/macOS)
\tazd exec "Write-Host 'Hello'" --shell pwsh   # Inline PowerShell
\tazd exec ./deploy.ps1 --shell pwsh            # Script with shell
\tazd exec ./build.sh -- --verbose              # Script with args
\tazd exec ./init.sh -i                         # Interactive mode`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Parse script arguments - everything after the script path
			scriptArgs := []string{}
			scriptInput := args[0]

			// Cobra doesn't parse args after -- automatically for us
			// They're in cmd.Flags().Args() after the script path
			if len(args) > 1 {
				scriptArgs = args[1:]
			}

			// Create executor
			exec := newScriptExecutor(executor.Config{
				Shell:               shell,
				Interactive:         interactive,
				StopOnKeyVaultError: stopOnKeyVaultError,
				Args:                scriptArgs,
			})

			// Check if input is a file or inline script
			// Try to resolve as file path first
			absPath, err := filepath.Abs(scriptInput)
			if err == nil {
				if _, statErr := os.Stat(absPath); statErr == nil {
					// It's a file that exists, execute as file
					return exec.Execute(cmd.Context(), absPath)
				}
			}

			// Not a file, treat as inline script
			return exec.ExecuteInline(cmd.Context(), scriptInput)
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Set output format from flag
			if outputFormat == "json" {
				if err := cliout.SetFormat("json"); err != nil {
					return fmt.Errorf("failed to set output format: %w", err)
				}
			}

			// Handle working directory change
			if cwd != "" {
				if err := os.Chdir(cwd); err != nil {
					return fmt.Errorf("failed to change working directory to %s: %w", cwd, err)
				}
			}

			// Handle debug mode
			if debugMode {
				_ = os.Setenv("AZD_DEBUG", "true")
			}

			// Handle no-prompt mode
			if noPrompt {
				_ = os.Setenv("AZD_NO_PROMPT", "true")
			}

			// Handle environment selection
			if environment != "" {
				// Load environment variables from the specified environment
				if err := env.LoadAzdEnvironment(cmd.Context(), environment); err != nil {
					return fmt.Errorf("failed to load environment '%s': %w", environment, err)
				}
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

	// Allow passthrough flags meant for the invoked command without requiring "--".
	rootCmd.FParseErrWhitelist.UnknownFlags = true
	// Stop flag parsing after the first script argument so downstream flags are preserved as args.
	rootCmd.Flags().SetInterspersed(false)
	rootCmd.PersistentFlags().SetInterspersed(false)

	// Add extension-specific flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Output format: default or json")

	// Add flags for direct script execution (when using 'azd exec ./script.sh')
	rootCmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run script in interactive mode")
	rootCmd.Flags().BoolVar(&stopOnKeyVaultError, "stop-on-keyvault-error", false, "Fail-fast: stop execution when any Key Vault reference fails to resolve")

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

	return rootCmd
}
