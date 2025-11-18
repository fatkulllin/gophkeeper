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

	cmds.AddCommand(NewCmdLogin())
	cmds.AddCommand(NewCmdRegister())

	return cmds
}
