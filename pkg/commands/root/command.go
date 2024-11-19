package root

import (
	"github.com/eddycharly/generic-auth-server/pkg/commands/serve"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	root := &cobra.Command{
		Use:   "generic-auth-server",
		Short: "generic-auth-server is a policy based authentication/authorization server",
	}
	root.AddCommand(serve.Command())
	return root
}
