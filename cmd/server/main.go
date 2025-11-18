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

// Пакет main запускает серверную часть GophKeeper.
// Приложение загружает конфигурацию, инициализирует логгер,
// создаёт экземпляр сервера и обеспечивает корректное завершение
// работы при получении системных сигналов.
func main() {

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	err = logger.Initialize(cfg.LogLevel, cfg.DevelopLog)

	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Log.Debug("Loaded config", zap.Any("config", cfg))

	application, err := app.NewApp(cfg)

	if err != nil {
		logger.Log.Fatal("failed to initialize gophkeeper", zap.Error(err))
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Log.Info("Starting gophkeeper server...")

	if err := application.Run(ctx); err != nil {
		logger.Log.Fatal("app shutdown with error", zap.Error(err))
	}

}
