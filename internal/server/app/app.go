package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/server/auth"
	"github.com/fatkulllin/gophkeeper/internal/server/config"
	"github.com/fatkulllin/gophkeeper/internal/server/cryptoutil"
	"github.com/fatkulllin/gophkeeper/internal/server/db"
	"github.com/fatkulllin/gophkeeper/internal/server/handlers"
	"github.com/fatkulllin/gophkeeper/internal/server/password"
	"github.com/fatkulllin/gophkeeper/internal/server/repositories/postgres"
	"github.com/fatkulllin/gophkeeper/internal/server/server"
	"github.com/fatkulllin/gophkeeper/internal/server/service"
	"github.com/fatkulllin/gophkeeper/migrations"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// App объединяет зависимости серверного приложения и
// управляет запуском HTTP и gRPC серверов.
type App struct {
	server *server.Server
	pgConn *sql.DB
}

// NewApp создаёт и настраивает серверное приложение GophKeeper.
// Здесь выполняется подключение к базе данных, применение миграций,
// инициализация сервисов, хендлеров и серверов.
//
// Возвращает экземпляр App или ошибку инициализации.
func NewApp(cfg config.Config) (App, error) {
	pgConn, err := db.NewPostgres(cfg.DatabaseURI)

	if err != nil {
		return App{}, fmt.Errorf("failed to connect to database %w", err)
	}

	logger.Log.Debug("successfully connected to database")

	err = db.Bootstrap(pgConn, migrations.FS)
	if err != nil {
		return App{}, fmt.Errorf("failed to apply migrations %w", err)
	}
	logger.Log.Debug("database migrated successfully")

	recordRepo := postgres.NewRecordRepo(pgConn)
	userRepo := postgres.NewUserRepo(pgConn)

	v := validator.New()
	tokenManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTExpires)

	logger.Log.Debug("init jwt manager successfully")

	pwdHasher := password.NewPassword()
	cryptoUtil := cryptoutil.NewCryptoUtil(cfg.MasterKey)

	service := service.NewService(userRepo, recordRepo, tokenManager, pwdHasher, cryptoUtil)
	healthHandler := handlers.NewHealthHandler()
	loggerHandler := handlers.NewLoggerHandler(v)
	authHandler := handlers.NewAuthHandler(service.User, v)
	recordHandler := handlers.NewRecordHandler(service.Record, v)
	srv := server.NewServer(cfg, healthHandler, loggerHandler, authHandler, recordHandler)

	return App{
		server: srv,
		pgConn: pgConn,
	}, nil
}

// Run запускает HTTP и gRPC сервера и ожидает их завершения.
// Остановка выполняется при получении сигнала завершения
// или при возникновении ошибки в одном из серверов.
func (app *App) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	// HTTP server
	group.Go(func() error {
		if err := app.server.Start(ctx); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("http server exited with error", zap.Error(err))
			return err
		}
		return nil
	})

	// gRPC server
	group.Go(func() error {
		if err := app.server.StartGRPC(ctx); err != nil {
			logger.Log.Error("grpc server exited with error", zap.Error(err))
			return err
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		logger.Log.Warn("shutting down due to error", zap.Error(err))
		return err
	}

	if err := app.pgConn.Close(); err != nil {
		logger.Log.Warn("failed to close database connection", zap.Error(err))
	}

	logger.Log.Info("shutting down...")

	logger.Log.Info("shutdown complete")

	return nil
}
