// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/jongio/azd-exec/cli/src/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command that displays extension version information.
func NewVersionCommand(outputFormat *string) *cobra.Command {
	return azdext.NewVersionCommand("jongio.azd.exec", version.Version, outputFormat)
}
