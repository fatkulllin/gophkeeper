/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

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

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	registerCmd.Flags().StringP("username", "u", "", "username")
	registerCmd.Flags().StringP("password", "p", "", "password")
}

// func registerUser(ctx context.Context, client *apiclient.Client, username, password, url string) (*models.Response, error) {

// 	user := map[string]string{
// 		"username": username,
// 		"password": password,
// 	}

// 	reqBody, err := json.Marshal(user)

// 	if err != nil {
// 		return &models.Response{}, fmt.Errorf("failed marshal batch: %w", err)
// 	}

// 	bodyReader := bytes.NewBuffer(reqBody)

// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)

// 	if err != nil {
// 		return &models.Response{}, fmt.Errorf("failed to create request: %w", err)
// 	}

// 	resp, err := client.Do(req)

// 	if err != nil {
// 		return &models.Response{}, fmt.Errorf("response error: %w", err)
// 	}

// 	return resp, nil
// }
