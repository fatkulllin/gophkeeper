package record

import "github.com/spf13/cobra"

func NewCmdRecord() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "record",
		Short: "Manager record",
	}
	cmds.AddCommand(NewCmdAdd())
	cmds.AddCommand(NewCmdGetAll())
	cmds.AddCommand(NewCmdGet())
	cmds.AddCommand(NewCmdDelete())
	cmds.AddCommand(NewCmdUpdate())
	cmds.AddCommand(NewCmdSync())
	return cmds
}
