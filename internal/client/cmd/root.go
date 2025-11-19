/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatkulllin/gophkeeper/internal/client/app"
	"github.com/fatkulllin/gophkeeper/internal/client/cmd/record"
	usermanager "github.com/fatkulllin/gophkeeper/internal/client/cmd/user"
	"github.com/fatkulllin/gophkeeper/internal/client/service"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	var svc *service.Service

	rootCmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "",
		Long:  ``,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			if err = initializeConfig(cmd); err != nil {
				return err
			}
			if err = initializeLogger(); err != nil {
				return err
			}

			svc, err = app.InitApp()
			if err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.PersistentFlags().String("log-level", "info", "logging level (debug, info, warn, error)")
	rootCmd.PersistentFlags().Bool("develop-log", false, "enable development logging")
	rootCmd.PersistentFlags().StringP("server", "s", "http://localhost:8080", "server address")
	rootCmd.AddCommand(usermanager.NewCmdUser(svc))
	rootCmd.AddCommand(record.NewCmdRecord(svc))
	rootCmd.AddCommand(NewCmdLogout(svc))
	return rootCmd
}

var (
	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initializeLogger() error {

	err := logger.Initialize(viper.GetString("log-level"), viper.GetBool("develop-log"))
	if err != nil {
		return err
	}
	return nil
}

func initializeConfig(cmd *cobra.Command) error {
	// 1. Set up Viper to use environment variables.
	viper.SetEnvPrefix("GOPHKEEPER")
	// Allow for nested keys in environment variables (e.g. `GOPHKEEPER_DATABASE_HOST`)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "_"))
	viper.AutomaticEnv()

	// 2. Handle the configuration file.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		// This is an optional but useful step to debug your config.
		fmt.Println("Configuration initialized. Using config file:", viper.ConfigFileUsed())
	}

	// 3. Read the configuration file.
	// If a config file is found, read it in. We use a robust error check
	// to ignore "file not found" errors, but panic on any other error.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. Bind Cobra flags to Viper.
	// This is the magic that makes the flag values available through Viper.
	// It binds the full flag set of the command passed in.
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		return err
	}

	return nil

}
