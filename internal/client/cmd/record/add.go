package record

import (
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdAdd() *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Create new record",
		RunE: func(cmd *cobra.Command, args []string) error {
			record := model.RecordInput{
				Type:     model.RecordType(viper.GetString("type")),
				Metadata: viper.GetString("metadata"),
				Data:     json.RawMessage(viper.GetString("data")),
			}
			url := viper.GetString("server") + "/api/record"
			resp, err := app.CliService.Record.Add(cmd.Context(), record, url)

			if err != nil {
				return fmt.Errorf("internal error: %v", err.Error())
			}
			if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
				if resp.StatusCode == 401 {
					return fmt.Errorf("unauthorized: %s", resp.Body)
				}
				return fmt.Errorf("registration failed: %s", string(resp.Body))
			}
			logger.Log.Info("record add successfully")

			return nil
		},
	}

	addCmd.Flags().String("type", "", "record type")
	addCmd.Flags().String("metadata", "", "record metadata")
	addCmd.Flags().String("data", "", "json with data")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("data")
	return addCmd
}
