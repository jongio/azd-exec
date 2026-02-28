// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
)

// NewListenCommand creates the listen command that starts the azd extension host.
// This command is invoked by azd to establish lifecycle event communication via gRPC.
func NewListenCommand() *cobra.Command {
	return azdext.NewListenCommand(nil)
}
