/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package usermanager

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the serve command
func NewCmdLogin(svc *service.Service) *cobra.Command {

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
			err := svc.User.ClearDB()
			if err != nil {
				return err
			}

			username := viper.GetString("username")
			password := viper.GetString("password")

			url := viper.GetString("server") + "/api/user/login" + "?userkey=true"

			ctx := context.Background()

			if username == "" || password == "" {
				return fmt.Errorf("username and password are required")
			}

			resp, err := svc.User.LoginUser(ctx, username, password, url)
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

			err = svc.User.SaveToken("token", auth_token)
			if err != nil {
				return fmt.Errorf("failed save token: %v", err)
			}

			var userKeyResponse model.UserKeyRespone
			err = json.Unmarshal(resp.Body, &userKeyResponse)
			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			err = svc.User.SaveUserKey(userKeyResponse.UserKey)
			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}

			urlRecords := viper.GetString("server") + "/api/records"
			recordsResponse, err := svc.Record.Get(cmd.Context(), urlRecords)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			if recordsResponse.StatusCode >= 400 && recordsResponse.StatusCode <= 500 {
				if recordsResponse.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", recordsResponse.Body)
				}
				return fmt.Errorf("registration failed: %s", string(recordsResponse.Body))
			}
			logger.Log.Info("get all successfully")
			var records []model.Record
			if err := json.Unmarshal(recordsResponse.Body, &records); err != nil {
				return fmt.Errorf("failed to parse JSON: %w", err)
			}
			if err := svc.Record.SaveRecords(records); err != nil {
				return fmt.Errorf("failed to save records to bolt: %w", err)
			}

			return nil
		},
	}
	cmd.Flags().StringP("username", "u", "", "username")
	cmd.Flags().StringP("password", "p", "", "password")
	cmd.Flags().Bool("userkey", false, "get user key")
	return cmd
}
