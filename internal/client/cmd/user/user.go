package usermanager

import (
	"context"

	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

func NewCmdUser(svc *service.Service, ctx context.Context) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "user",
		Short: "Manager user",
	}

	cmds.AddCommand(NewCmdLogin(svc))
	cmds.AddCommand(NewCmdRegister(svc))

	return cmds
}
