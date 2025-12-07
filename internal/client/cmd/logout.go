/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

// loginCmd represents the serve command
func NewCmdLogout(svc *service.Service) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) error {

			err := svc.User.ClearDB()
			if err != nil {
				return err
			}

			err = svc.User.ClearToken("token")

			if err != nil {
				return err
			}

			return nil
		},
	}
	return cmd
}
