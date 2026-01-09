// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"github.com/spf13/cobra"
)

// NewListenCommand creates a new listen command for the azd extension framework.
// This command is required by the azd extension framework for inter-process communication.
// It is marked as hidden since it's an internal implementation detail not meant for direct use.
// The extension currently operates in "exec mode" without a persistent listener,
// so this command is a no-op placeholder for framework compatibility.
func NewListenCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "listen",
		Short:  "Start extension listener (internal use only)",
		Long:   `Start the extension listener for the azd extension framework. This command is used internally by azd and should not be called directly.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// This is a placeholder for the extension framework's listen functionality.
			// In a full implementation, this would start a gRPC server for azd communication.
			// For now, this extension operates in "exec mode" without persistent listener,
			// so returning nil is appropriate.
			return nil
		},
	}
}
