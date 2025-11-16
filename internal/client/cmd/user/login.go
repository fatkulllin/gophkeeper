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
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the serve command
func NewCmdLogin() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate an existing user",
		Long: `Authenticate an existing user with the GophKeeper server.

Examples:
  gophkeeper login -u alice -p secret123
  gophkeeper login --username bob --password mypass

After successful authentication, your access token is stored locally
and used for future requests.`,
		RunE: func(cmd *cobra.Command, args []string) error {

			username := viper.GetString("username")
			password := viper.GetString("password")
			ctx := context.Background()
			url := viper.GetString("server") + "/api/user/login"

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

			logger.Log.Info("loggin successfully")

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
	cmd.Flags().StringP("username", "u", "", "username")
	cmd.Flags().StringP("password", "p", "", "password")
	return cmd
}
