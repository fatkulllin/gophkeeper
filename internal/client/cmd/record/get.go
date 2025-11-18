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

func NewCmdGet() *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get record",
		RunE: func(cmd *cobra.Command, args []string) error {
			remote := viper.GetBool("remote")

			if remote {
				url := viper.GetString("server") + "/api/records/" + viper.GetString("id")
				resp, err := app.CliService.Record.Get(cmd.Context(), url)

				if err != nil {
					return fmt.Errorf("internal error: %v", err.Error())
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

				logger.Log.Info("get record successfully")

				var pretty bytes.Buffer
				err = json.Indent(&pretty, resp.Body, "", "  ")
				if err != nil {
					fmt.Println(string(resp.Body))
				} else {
					fmt.Println(pretty.String())
				}

				return nil
			}
			id := viper.GetInt64("id")
			record, err := app.CliService.Record.GetLocal(cmd.Context(), id)

			if err != nil {
				logger.Log.Error("", zap.Error(err))
				return fmt.Errorf("internal error: %v", err.Error())
			}

			rawJSON, err := json.Marshal(record)
			if err != nil {
				logger.Log.Error("", zap.Error(err))
				return fmt.Errorf("internal error: %v", err.Error())
			}

			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, rawJSON, "", "  ")
			if err != nil {
				fmt.Println(record)
			} else {
				fmt.Println(prettyJSON.String())
			}

			return nil
		},
	}
	getCmd.Flags().String("id", "", "id record")
	getCmd.MarkFlagRequired("id")
	getCmd.Flags().Bool("remote", false, "fetch records from server instead of local bbolt")
	return getCmd
}
