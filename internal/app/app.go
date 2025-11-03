package app

import (
	"context"
	"net/http"
	"sync"

	"github.com/fatkulllin/gophkeeper/internal/config"
	"github.com/fatkulllin/gophkeeper/internal/handlers"
	"github.com/fatkulllin/gophkeeper/internal/logger"
	"github.com/fatkulllin/gophkeeper/internal/server"
	"go.uber.org/zap"
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
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	var wg sync.WaitGroup

	errCh := make(chan error, 2)

	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := app.server.Start(ctx); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server exited with error", zap.Error(err))
			errCh <- err
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		logger.Log.Info("shutting down...")
	case err := <-errCh:
		logger.Log.Warn("shutting down due to error")
		return err
	}

	wg.Wait()
	logger.Log.Info("shutdown complete")
	return nil
}
