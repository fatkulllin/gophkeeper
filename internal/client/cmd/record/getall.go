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

func NewCmdGetAll() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "getall",
		Short: "Get all records",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := viper.GetString("server") + "/api/records"
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
			logger.Log.Info("get all successfully")
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
	return addCmd
}
