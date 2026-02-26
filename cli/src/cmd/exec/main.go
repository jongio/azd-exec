// Package main provides the azd exec extension command-line interface.
// It enables execution of scripts with full Azure Developer CLI context,
// including environment variables and Azure credentials.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/jongio/azd-core/cliout"
	"github.com/jongio/azd-core/env"
	"github.com/jongio/azd-exec/cli/src/cmd/exec/commands"
	"github.com/jongio/azd-exec/cli/src/internal/executor"
	"github.com/jongio/azd-exec/cli/src/internal/skills"
	"github.com/jongio/azd-exec/cli/src/internal/version"
	"github.com/spf13/cobra"
)

var (
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
	rootCmd, extCtx := azdext.NewExtensionRootCommand(azdext.ExtensionCommandOptions{
		Name:    "exec",
		Version: version.Version,
		Use:     "exec [script-file-or-command] [-- script-args...]",
		Short:   "Exec - Execute commands/scripts with Azure Developer CLI context",
		Long: `Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.

Examples:
\tazd exec ./setup.sh                           # Execute script file
\tazd exec 'echo $AZURE_ENV_NAME'               # Inline bash (Linux/macOS)
\tazd exec --shell pwsh "Write-Host 'Hello'"   # Inline PowerShell
\tazd exec --shell pwsh ./deploy.ps1            # Script with shell
\tazd exec ./build.sh -- --verbose              # Script with args
\tazd exec ./init.sh -i                         # Interactive mode`,
	})

	rootCmd.Args = cobra.MinimumNArgs(1)
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
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
	}

	// Save the SDK's PersistentPreRunE so we can chain it
	sdkPreRunE := rootCmd.PersistentPreRunE
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// Run SDK setup first (trace context, cwd, debug, etc.)
		if sdkPreRunE != nil {
			if err := sdkPreRunE(cmd, args); err != nil {
				return err
			}
		}

		// Set output format from flag
		if extCtx.OutputFormat == "json" {
			if err := cliout.SetFormat("json"); err != nil {
				return fmt.Errorf("failed to set output format: %w", err)
			}
		}

		// Handle environment selection
		if extCtx.Environment != "" {
			// Load environment variables from the specified environment
			if err := env.LoadAzdEnvironment(cmd.Context(), extCtx.Environment); err != nil {
				return fmt.Errorf("failed to load environment '%s': %w", extCtx.Environment, err)
			}
		}

		// Install Copilot skill
		if err := skills.InstallSkill(); err != nil {
			if extCtx.Debug {
				fmt.Fprintf(os.Stderr, "Warning: failed to install copilot skill: %v\n", err)
			}
		}

		return nil
	}

	// Allow passthrough flags meant for the invoked command without requiring "--".
	rootCmd.FParseErrWhitelist.UnknownFlags = true
	// Stop flag parsing after the first script argument so downstream flags are preserved as args.
	rootCmd.Flags().SetInterspersed(false)
	rootCmd.PersistentFlags().SetInterspersed(false)

	// Add flags for direct script execution (when using 'azd exec ./script.sh')
	rootCmd.Flags().StringVarP(&shell, "shell", "s", "", "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.")
	rootCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run script in interactive mode")
	rootCmd.Flags().BoolVar(&stopOnKeyVaultError, "stop-on-keyvault-error", false, "Fail-fast: stop execution when any Key Vault reference fails to resolve")

	// Register subcommands
	rootCmd.AddCommand(
		commands.NewVersionCommand(&extCtx.OutputFormat),
		commands.NewListenCommand(),
		commands.NewMetadataCommand(newRootCmd),
		commands.NewMCPCommand(),
	)

	return rootCmd
}
