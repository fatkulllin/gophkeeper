/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package usermanager

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/internal/client/filemanager"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdRegister() *cobra.Command {

	// registerCmd represents the register command
	var registerCmd = &cobra.Command{
		Use:   "register",
		Short: "Create a new user account",
		Long: `Register a new user account on the GophKeeper server.

Examples:
  gophkeeper register -u newuser -p strongpass123

Upon successful registration, you will be automatically logged in.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			username := viper.GetString("username")
			password := viper.GetString("password")
			ctx := context.Background()
			url := viper.GetString("server") + "/api/user/register"

			if username == "" || password == "" {
				return fmt.Errorf("username and password are required")
			}

			resp, err := app.CliService.User.LoginUser(ctx, username, password, url)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}

			if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
				if resp.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", resp.Body)
				}
				return fmt.Errorf("registration failed: %s", string(resp.Body))
			}
			logger.Log.Info("registration successfully")

			var auth_token string

			for _, c := range resp.Cookies {

				if c.Name == "auth_token" {
					auth_token = c.Value
					break
				}
			}
			permission, err := strconv.ParseUint("0600", 8, 32)

			if err != nil {
				return err
			}

			filemanager.NewFileManager().SaveFile("token", auth_token, os.FileMode(permission))
			return nil
		},
	}
	registerCmd.Flags().StringP("username", "u", "", "username")
	registerCmd.Flags().StringP("password", "p", "", "password")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
	return registerCmd
}
