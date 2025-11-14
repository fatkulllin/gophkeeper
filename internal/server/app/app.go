package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/server/auth"
	"github.com/fatkulllin/gophkeeper/internal/server/config"
	"github.com/fatkulllin/gophkeeper/internal/server/handlers"
	"github.com/fatkulllin/gophkeeper/internal/server/password"
	"github.com/fatkulllin/gophkeeper/internal/server/repository/postgres"
	"github.com/fatkulllin/gophkeeper/internal/server/server"
	"github.com/fatkulllin/gophkeeper/internal/server/service"
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/fatkulllin/gophkeeper/migrations"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	server *server.Server
	pgRepo *postgres.PGRepo
}

func NewApp(cfg config.Config) (App, error) {

	pgRepo, err := postgres.NewPGRepo(cfg.DatabaseURI)

	if err != nil {
		return App{}, fmt.Errorf("connect to Database is unavailable: %w", err)
	}

	logger.Log.Debug("successfully connected to database")

	err = pgRepo.Bootstrap(migrations.FS)

	if err != nil {
		return App{}, fmt.Errorf("migrate is not run: %w", err)
	}

	logger.Log.Debug("database migrated successfully")

	tokenManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpires)

	logger.Log.Debug("init jwt manager successfully")

	password := password.NewPassword()

	service := service.NewService(pgRepo, tokenManager, password)
	healthHandler := handlers.NewHealthHandler()
	loggerHandler := handlers.NewLoggerHandler()
	authHandler := handlers.NewAuthHandler(service)
	recordHandler := handlers.NewRecordHandler(service)
	server := server.NewServer(cfg, healthHandler, loggerHandler, authHandler, recordHandler)

	return App{
		server: server,
		pgRepo: pgRepo,
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
