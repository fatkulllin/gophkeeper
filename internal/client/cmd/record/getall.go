package record

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewCmdGetAll() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "getall",
		Short: "Get all records",
		RunE: func(cmd *cobra.Command, args []string) error {
			remote := viper.GetBool("remote")

			if remote {
				url := viper.GetString("server") + "/api/records"
				resp, err := app.CliService.Record.Get(cmd.Context(), url)

				if err != nil {
					return fmt.Errorf("internal error: %w", err)
				}

				if resp.StatusCode >= 400 {
					if resp.StatusCode == 401 {
						return fmt.Errorf("unauthorized: %s", resp.Body)
					}

					return fmt.Errorf(
						"failed to fetch records: status=%d body=%s",
						resp.StatusCode,
						resp.Body,
					)
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
			}
			records, err := app.CliService.Record.GetAll()
			if err != nil {
				logger.Log.Error("", zap.Error(err))
				return fmt.Errorf("internal error: %v", err.Error())
			}
			rawJSON, err := json.Marshal(records)
			if err != nil {
				logger.Log.Error("", zap.Error(err))
				return fmt.Errorf("internal error: %v", err.Error())
			}

			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, rawJSON, "", "  ")
			if err != nil {
				fmt.Println(records)
			} else {
				fmt.Println(prettyJSON.String())
			}
			return nil
		},
	}
	getCmd.Flags().Bool("remote", false, "fetch records from server instead of local bbolt")
	return getCmd
}
