package record

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdUpdate() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "update",
		Short: "update record",
		RunE: func(cmd *cobra.Command, args []string) error {
			metadata := viper.GetString("metadata")
			data := json.RawMessage(viper.GetString("data"))
			record := model.RecordUpdateInput{}
			if !cmd.Flags().Changed("metadata") && !cmd.Flags().Changed("data") {
				return fmt.Errorf("metadata and data are required")
			}
			if cmd.Flags().Changed("metadata") {
				record.Metadata = &metadata
			}
			if cmd.Flags().Changed("data") {
				record.Data = &data
			}

			url := viper.GetString("server") + "/api/records/" + viper.GetString("id")
			resp, err := app.CliService.Record.Update(cmd.Context(), url, record)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
				if resp.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", resp.Body)
				}
				return fmt.Errorf("registration failed: %s", string(resp.Body))
			}
			logger.Log.Info("update record successfully")
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
	addCmd.Flags().String("metadata", "", "metadata record")
	addCmd.Flags().String("data", "", "data record")
	addCmd.Flags().String("id", "", "id record")
	addCmd.MarkFlagRequired("id")
	return addCmd
}
