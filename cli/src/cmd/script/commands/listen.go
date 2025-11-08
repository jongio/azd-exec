package commands

import (
	"github.com/spf13/cobra"
)

// NewListenCommand creates a new listen command.
// This is required by the azd extension framework for inter-process communication.
func NewListenCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "listen",
		Short:  "Start extension listener (internal use only)",
		Long:   `Start the extension listener for the azd extension framework. This command is used internally by azd and should not be called directly.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This is a placeholder for the extension framework's listen functionality
			// In a full implementation, this would start a gRPC server for azd communication
			return nil
		},
	}
}
