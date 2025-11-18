/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/spf13/cobra"
)

// loginCmd represents the serve command
func NewCmdLogout() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := app.CliService.User.ClearDB()
			if err != nil {
				return err
			}

			err = app.CliService.User.ClearToken("token")

			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
