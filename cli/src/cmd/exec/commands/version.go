// Package commands provides subcommands for the azd exec extension.
package commands

import (
	coreversion "github.com/jongio/azd-core/version"
	"github.com/jongio/azd-exec/cli/src/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand creates a new version command that displays extension version information.
func NewVersionCommand(outputFormat *string) *cobra.Command {
	return coreversion.NewCommand(version.Info, outputFormat)
}
