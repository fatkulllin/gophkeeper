package record

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewCmdDelete() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete record",
		RunE: func(cmd *cobra.Command, args []string) error {
			idRecord := viper.GetString("id")
			url := viper.GetString("server") + "/api/records/" + idRecord
			resp, err := app.CliService.Record.Delete(cmd.Context(), url)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
				if resp.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", resp.Body)
				}
				return fmt.Errorf("registration failed: %s", string(resp.Body))
			}
			logger.Log.Info("delete record successfully", zap.String("record id", idRecord))
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
