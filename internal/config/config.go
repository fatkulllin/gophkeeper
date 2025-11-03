package config

import (
	"fmt"
	"net"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	HTTPAddress string `env:"HTTP_ADDRESS" envDefault:"localhost:8080"`
	GRPCAddress string `env:"GRPC_ADDRESS" envDefault:"localhost:9090"`
	DevelopLog  bool   `env:"DEVELOP_LOG" envDefault:"false"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func validateAddress(s string) error {
	_, _, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (Config, error) {

	config := Config{}

	pflag.CommandLine.SortFlags = false // чтобы флаги выводились в заданном порядке
	pflag.StringVar(&config.HTTPAddress, "http-address", config.HTTPAddress, "HTTP server listen address (host:port)")
	pflag.StringVar(&config.GRPCAddress, "grpc-address", config.GRPCAddress, "GRPC server listen address (host:port)")
	pflag.StringVarP(&config.LogLevel, "log-level", "l", config.LogLevel, "logging level: debug, info, warn, error")
	pflag.BoolVarP(&config.DevelopLog, "develop-log", "d", config.DevelopLog, "enabled develop log")

	pflag.Parse()

	err := env.Parse(&config)

	if err != nil {
		return config, fmt.Errorf("error parsing environment %w", err)
	}

	if err := validateAddress(config.HTTPAddress); err != nil {
		return config, fmt.Errorf("invalid server address: %s, %w", config.HTTPAddress, err)
	}

	return config, nil
}
