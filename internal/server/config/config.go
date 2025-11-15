package config

import (
	"fmt"
	"net"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct {
	HTTPAddress string `env:"HTTP_ADDRESS"`
	GRPCAddress string `env:"GRPC_ADDRESS"`
	DevelopLog  bool   `env:"DEVELOP_LOG"`
	LogLevel    string `env:"LOG_LEVEL"`
	DatabaseURI string `env:"DATABASE_URI"`
	JWTSecret   string `env:"JWT_SECRET_KEY"`
	JWTExpires  int    `env:"JWT_EXPIRES"`
	MasterKey   string `env:"MASTER_KEY"`
}

const (
	DefaultHTTPAddress = "localhost:8080"
	DeafultGRPCAddress = "localhost:9090"
	DeafultDevelopLog  = false
	DefaultLogLevel    = "INFO"
	DefaultDatabaseURI = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	DefaultJWTSecret   = "TOKEN"
	DefaultJWTExpires  = 24
	DefaultMasterKey   = "M2WXo9hmd9pHMoaY8UHQokJhGXetCNBP"
)

func validateAddress(s string) error {
	_, _, err := net.SplitHostPort(s)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() (Config, error) {

	config := Config{
		HTTPAddress: DefaultHTTPAddress,
		GRPCAddress: DeafultGRPCAddress,
		DevelopLog:  DeafultDevelopLog,
		LogLevel:    DefaultLogLevel,
		DatabaseURI: DefaultDatabaseURI,
		JWTSecret:   DefaultJWTSecret,
		JWTExpires:  DefaultJWTExpires,
		MasterKey:   DefaultMasterKey,
	}

	pflag.CommandLine.SortFlags = false // чтобы флаги выводились в заданном порядке
	pflag.StringVar(&config.HTTPAddress, "http-address", config.HTTPAddress, "HTTP server listen address (host:port)")
	pflag.StringVar(&config.GRPCAddress, "grpc-address", config.GRPCAddress, "GRPC server listen address (host:port)")
	pflag.StringVarP(&config.LogLevel, "log-level", "l", config.LogLevel, "logging level: debug, info, warn, error")
	pflag.BoolVar(&config.DevelopLog, "develop-log", config.DevelopLog, "enabled develop log")
	pflag.StringVarP(&config.DatabaseURI, "database", "d", config.DatabaseURI, "set database dsn")
	pflag.StringVarP(&config.JWTSecret, "secret", "s", config.JWTSecret, "set secret token")
	pflag.IntVarP(&config.JWTExpires, "expires", "e", config.JWTExpires, "set expires jwt")
	pflag.StringVarP(&config.MasterKey, "master-key", "m", config.MasterKey, "set master key")
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
