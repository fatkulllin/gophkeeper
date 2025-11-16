package usermanager

import (
	"github.com/spf13/cobra"
)

func NewCmdUser() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "user",
		Short: "Manager user",
	}

	// cmds.AddCommand(NewCmdCanI(f, streams))
	// cmds.AddCommand(NewCmdReconcile(f, streams))
	// cmds.AddCommand(NewCmdWhoAmI(f, streams))
	cmds.AddCommand(NewCmdLogin())
	cmds.AddCommand(NewCmdRegister())

	return cmds
}
