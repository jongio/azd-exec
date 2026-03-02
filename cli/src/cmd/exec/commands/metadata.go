// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
)

// NewMetadataCommand creates a metadata command that generates extension metadata
// using the official azdext SDK metadata generator.
func NewMetadataCommand(rootCmdProvider func() *cobra.Command) *cobra.Command {
	return azdext.NewMetadataCommand("1.0", "jongio.azd.exec", rootCmdProvider)
}
