package app

import (
	"context"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/config"
	"github.com/fatkulllin/gophkeeper/internal/handlers"
	"github.com/fatkulllin/gophkeeper/internal/logger"
	"github.com/fatkulllin/gophkeeper/internal/server"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	server server.Server
}

func NewApp(cfg config.Config) (App, error) {

	healthHandler := handlers.NewHealthHandler()
	loggerHandler := handlers.NewLoggerHandler()
	server := server.NewServer(cfg, healthHandler, loggerHandler)

	return App{
		server: server,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		if err := app.server.Start(ctx); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server exited with error", zap.Error(err))
			return err
		}
		return nil
	})

	group.Go(func() error {
		if err := app.server.StartGRPC(ctx); err != nil {
			logger.Log.Error("server exited with error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		logger.Log.Warn("shutting down due to error", zap.Error(err))
		return err
	}

	logger.Log.Info("shutting down...")

	logger.Log.Info("shutdown complete")
	return nil
}
