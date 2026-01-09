// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"encoding/json"
	"fmt"
	"os"

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
			switch *outputFormat {
			case "json":
				output := map[string]string{
					"version": version.Version,
				}
				data, err := json.MarshalIndent(output, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting JSON: %v\n", err)
					return
				}
				fmt.Println(string(data))
			default:
				if quiet {
					fmt.Println(version.Version)
				} else {
					fmt.Printf("azd exec version %s\n", version.Version)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Display only the version number")
	return cmd
}
