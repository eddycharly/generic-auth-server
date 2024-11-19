package serve

import (
	authzserver "github.com/eddycharly/generic-auth-server/pkg/commands/serve/authz-server"
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run policy based authentication/authorization servers",
	}
	command.AddCommand(authzserver.Command())
	return command
}
