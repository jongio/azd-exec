// Package commands provides subcommands for the azd exec extension.
package commands

import (
	"context"
	"fmt"

	"github.com/azure/azure-dev/cli/azd/pkg/azdext"
	"github.com/spf13/cobra"
)

// NewListenCommand creates the listen command that starts the azd extension host.
// This command is invoked by azd to establish lifecycle event communication via gRPC.
func NewListenCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "listen",
		Short:  "Start extension listener (internal use only)",
		Long:   `Start the extension listener for the azd extension framework. This command is used internally by azd and should not be called directly.`,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := azdext.WithAccessToken(cmd.Context())

			azdClient, err := azdext.NewAzdClient()
			if err != nil {
				return fmt.Errorf("failed to create azd client: %w", err)
			}
			defer azdClient.Close()

			host := azdext.NewExtensionHost(azdClient).
				WithProjectEventHandler("postprovision", handlePostProvision).
				WithProjectEventHandler("postdeploy", handlePostDeploy).
				WithServiceEventHandler("postdeploy", handleServicePostDeploy, nil)

			if err := host.Run(ctx); err != nil {
				return fmt.Errorf("failed to run extension: %w", err)
			}

			return nil
		},
	}
}

func handlePostProvision(ctx context.Context, args *azdext.ProjectEventArgs) error {
	fmt.Printf("Post-provision completed for project: %s\n", args.Project.Name)
	return nil
}

func handlePostDeploy(ctx context.Context, args *azdext.ProjectEventArgs) error {
	fmt.Printf("Deployment completed for project: %s\n", args.Project.Name)
	return nil
}

func handleServicePostDeploy(ctx context.Context, args *azdext.ServiceEventArgs) error {
	fmt.Printf("Service %s deployed successfully\n", args.Service.Name)
	for _, artifact := range args.ServiceContext.Deploy {
		if artifact.Kind == azdext.ArtifactKind_ARTIFACT_KIND_ENDPOINT {
			fmt.Printf("  Endpoint: %s\n", artifact.Location)
		}
	}
	return nil
}
