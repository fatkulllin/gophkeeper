package record

import (
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdSync(svc *service.Service) *cobra.Command {
	addCmd := &cobra.Command{
		Use:   "sync",
		Short: "sync all records",
		RunE: func(cmd *cobra.Command, args []string) error {
			url := viper.GetString("server") + "/api/records"
			resp, err := svc.Record.Get(cmd.Context(), url)

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
			var records []model.Record
			if err := json.Unmarshal(resp.Body, &records); err != nil {
				return fmt.Errorf("failed to parse JSON: %w", err)
			}
			if err := svc.Record.SaveRecords(records); err != nil {
				return fmt.Errorf("failed to save records to bolt: %w", err)
			}
			return nil
		},
	}
	return addCmd
}
