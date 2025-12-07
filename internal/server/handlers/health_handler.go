// Package handlers содержит реализацию HTTP-хендлеров серверного приложения.
//
// HealthHandler предоставляет простой эндпоинт для проверки доступности сервиса.
package handlers

import (
	"net/http"

	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// HealthHandler обрабатывает запросы проверки состояния сервера.
// Используется системой мониторинга или балансировщиком для проверки,
// что HTTP-сервер запущен и отвечает на запросы.
type HealthHandler struct{}

// NewHealthHandler создаёт новый HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthHTTP возвращает HTTP 200 OK.
// Эндпоинт не выполняет дополнительных проверок и служит индикатором того,
// что HTTP-сервер доступен и готов обрабатывать запросы.
//
// GET /healthcheck
func (h *HealthHandler) HealthHTTP(res http.ResponseWriter, req *http.Request) {
	body := []byte("OK")
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(body); err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}
