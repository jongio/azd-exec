package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

// NewVersionCommand creates a new version command.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the extension version",
		Long:  `Display the version information for the azd exec extension.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("azd exec (azd-script) version %s\n", version)
		},
	}
}
