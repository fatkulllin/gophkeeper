package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatkulllin/gophkeeper/internal/auth"
	"github.com/fatkulllin/gophkeeper/internal/config"
	"github.com/fatkulllin/gophkeeper/internal/handlers"
	"github.com/fatkulllin/gophkeeper/internal/logger"
	logging "github.com/fatkulllin/gophkeeper/internal/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Server struct {
	config     config.Config
	httpServer *http.Server
}

// NewRouter создаёт и настраивает HTTP-роутер с хендлерами и middleware.
// Использует chi.Router и возвращает готовый маршрутизатор.
func NewRouter(healthHandler *handlers.HealthHandler, loggerHandler *handlers.LoggerHandler, authHandler *handlers.AuthHandler, jwtSecret string) chi.Router {
	r := chi.NewRouter()
	r.Use(logging.RequestLogger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Get("/healthcheck", healthHandler.HealthHTTP)
	r.Get("/debug/loglevel", loggerHandler.GetLevel)
	r.Post("/debug/loglevel", loggerHandler.SetLevel)
	r.Post("/api/user/register", authHandler.UserRegister)
	r.Post("/api/user/login", authHandler.UserLogin)
	r.Post("/api/user/logout", authHandler.UserLogout)
	r.Use(auth.AuthMiddleware(jwtSecret))
	return r
}

func NewServer(cfg config.Config, debugHandler *handlers.HealthHandler, loggerHandler *handlers.LoggerHandler, authHandler *handlers.AuthHandler) *Server {
	router := NewRouter(debugHandler, loggerHandler, authHandler, cfg.JWTSecret)
	return &Server{
		config: cfg,
		httpServer: &http.Server{
			Addr:         cfg.HTTPAddress,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (server *Server) Start(ctx context.Context) error {

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		logger.Log.Info("HTTP server shutting down...")

		if err := server.httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Log.Error("HTTP server shutdown failed", zap.Error(err))
		}
	}()

	logger.Log.Info("HTTP server started on", zap.String("address", server.httpServer.Addr))

	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP server error: %w", err)
	}

	logger.Log.Info("HTTP server stopped gracefully")
	return nil

}
