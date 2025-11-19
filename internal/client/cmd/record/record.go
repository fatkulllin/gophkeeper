package record

import (
	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

func NewCmdRecord(svc *service.Service) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "record",
		Short: "Manager record",
	}
	cmds.AddCommand(NewCmdAdd(svc))
	cmds.AddCommand(NewCmdGetAll(svc))
	cmds.AddCommand(NewCmdGet(svc))
	cmds.AddCommand(NewCmdDelete(svc))
	cmds.AddCommand(NewCmdUpdate(svc))
	cmds.AddCommand(NewCmdSync(svc))
	return cmds
}
