package record

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdGet() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "get",
		Short: "Get record",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := viper.GetString("server") + "/api/records/" + viper.GetString("id")
			resp, err := app.CliService.Record.Get(cmd.Context(), url)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
				if resp.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", resp.Body)
				}
				return fmt.Errorf("registration failed: %s", string(resp.Body))
			}
			logger.Log.Info("get record successfully")
			var pretty bytes.Buffer
			err = json.Indent(&pretty, resp.Body, "", "  ")
			if err != nil {
				fmt.Println(string(resp.Body))
			} else {
				fmt.Println(pretty.String())
			}

			return nil
		},
	}
	addCmd.Flags().String("id", "", "id record")
	addCmd.MarkFlagRequired("id")
	return addCmd
}
