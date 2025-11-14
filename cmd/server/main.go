package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatkulllin/gophkeeper/internal/server/app"
	"github.com/fatkulllin/gophkeeper/internal/server/config"
	"github.com/fatkulllin/gophkeeper/logger"
	"go.uber.org/zap"
)

func main() {

	config, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	err = logger.Initialize(config.LogLevel, config.DevelopLog)

	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Log.Debug("Loaded config", zap.Any("config", config))

	app, err := app.NewApp(config)

	if err != nil {
		logger.Log.Fatal("failed to initialize gophermart: ", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		logger.Log.Fatal("app shutdown with error", zap.Error(err))
	}

}
