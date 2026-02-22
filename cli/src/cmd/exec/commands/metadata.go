// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type metadataOutput struct {
	SchemaVersion string            `json:"schemaVersion"`
	ID            string            `json:"id"`
	Commands      []commandMetadata `json:"commands"`
}

type commandMetadata struct {
	Name        []string          `json:"name"`
	Short       string            `json:"short"`
	Long        string            `json:"long,omitempty"`
	Usage       string            `json:"usage,omitempty"`
	Examples    []exampleMetadata `json:"examples,omitempty"`
	Flags       []flagMetadata    `json:"flags,omitempty"`
	Subcommands []commandMetadata `json:"subcommands,omitempty"`
}

type exampleMetadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

type flagMetadata struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Default     string `json:"default,omitempty"`
}

// NewMetadataCommand creates a hidden command that outputs extension metadata as JSON.
func NewMetadataCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "metadata",
		Short:  "Generate extension metadata",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			metadata := metadataOutput{
				SchemaVersion: "1.0",
				ID:            "jongio.azd.exec",
				Commands: []commandMetadata{
					{
						Name:  []string{"exec"},
						Short: "Exec - Execute commands/scripts with Azure Developer CLI context",
						Long:  "Exec is an Azure Developer CLI extension that executes commands and scripts with full access to azd environment variables and configuration.",
						Usage: "exec [script-file-or-command] [-- script-args...]",
						Examples: []exampleMetadata{
							{Name: "execute script", Description: "Execute a script file with azd context", Usage: "azd exec ./my-script.sh"},
							{Name: "with shell", Description: "Execute with specific shell", Usage: "azd exec --shell pwsh ./deploy.ps1"},
							{Name: "inline command", Description: "Execute an inline command", Usage: "azd exec \"echo hello\""},
							{Name: "version", Description: "Display extension version", Usage: "azd exec version"},
						},
						Flags: []flagMetadata{
							{Name: "shell", Shorthand: "s", Description: "Shell to use for execution (bash, sh, zsh, pwsh, powershell, cmd). Auto-detected if not specified.", Type: "string"},
							{Name: "interactive", Shorthand: "i", Description: "Run script in interactive mode", Type: "bool"},
							{Name: "stop-on-keyvault-error", Description: "Fail-fast: stop execution when any Key Vault reference fails to resolve", Type: "bool"},
							{Name: "output", Shorthand: "o", Description: "Output format: default or json", Type: "string", Default: "default"},
							{Name: "debug", Description: "Enable debug mode", Type: "bool"},
							{Name: "no-prompt", Description: "Disable prompts", Type: "bool"},
							{Name: "cwd", Shorthand: "C", Description: "Sets the current working directory", Type: "string"},
							{Name: "environment", Shorthand: "e", Description: "The name of the environment to use", Type: "string"},
						},
						Subcommands: []commandMetadata{
							{
								Name:  []string{"exec", "version"},
								Short: "Display extension version",
							},
						},
					},
				},
			}

			out, err := json.MarshalIndent(metadata, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal metadata: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(out))
			return nil
		},
	}
}
