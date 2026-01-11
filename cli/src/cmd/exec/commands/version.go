// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"encoding/json"
	"fmt"

	"github.com/jongio/azd-core/cliout"
	"github.com/jongio/azd-exec/cli/src/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command that displays extension version information.
// The command supports both human-readable and JSON output formats.
// The outputFormat parameter controls the output style via the --output/-o flag.
func NewVersionCommand(outputFormat *string) *cobra.Command {
	var quiet bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display the extension version",
		Long:  `Display the version information for the azd exec extension.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Set output format from flag
			if *outputFormat == "json" {
				if err := cliout.SetFormat("json"); err != nil {
					cliout.Error("Failed to set output format: %v", err)
					return
				}
			} else {
				if err := cliout.SetFormat("default"); err != nil {
					cliout.Error("Failed to set output format: %v", err)
					return
				}
			}

			if cliout.IsJSON() {
				// JSON output mode
				output := map[string]string{
					"version": version.Version,
				}
				data, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					cliout.Error("Error formatting JSON: %v", err)
					return
				}
				fmt.Println(string(data))
			} else {
				// Human-readable output with colors
				if quiet {
					fmt.Println(version.Version)
				} else {
					cliout.Header("azd exec")
					cliout.Label("Version", version.Version)
					if version.BuildDate != "unknown" {
						cliout.Label("Build Date", version.BuildDate)
					}
					if version.GitCommit != "unknown" {
						cliout.Label("Git Commit", version.GitCommit)
					}
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Display only the version number")
	return cmd
}
