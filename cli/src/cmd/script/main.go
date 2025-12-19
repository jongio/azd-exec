package main

import (
	"fmt"
	"os"

	"github.com/jongio/azd-script/cli/src/cmd/script/commands"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	debugMode    bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "exec",
		Short: "Exec - Execute commands/scripts with Azure Developer CLI context",
		Long:  `Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if debugMode {
				_ = os.Setenv("AZD_SCRIPT_DEBUG", "true")
			}
			return nil
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "default", "Output format (default, json)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug logging")

	// Register all commands
	rootCmd.AddCommand(
		commands.NewRunCommand(),
		commands.NewVersionCommand(),
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
