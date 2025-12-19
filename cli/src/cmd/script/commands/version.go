package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

// NewVersionCommand creates a new version command.
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
					"version": version,
				}
				data, _ := json.MarshalIndent(output, "", "  ")
				fmt.Println(string(data))
			default:
				if quiet {
					fmt.Println(version)
				} else {
					fmt.Printf("azd exec version %s\n", version)
				}
			}
		},
	}

	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Display only the version number")
	return cmd
}
