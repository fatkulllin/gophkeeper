// Package handlers содержит реализацию HTTP-хендлеров приложения.
// LoggerHandler предоставляет эндпоинты для получения и изменения
// уровня логирования во время работы сервера.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// LoggerHandler обрабатывает запросы, связанные с управлением уровнем логирования.
// Позволяет получить текущий уровень логов и изменить его через HTTP API.
type LoggerHandler struct {
	validate *validator.Validate
}

// NewLoggerHandler создаёт новый LoggerHandler.
func NewLoggerHandler(validate *validator.Validate) *LoggerHandler {
	return &LoggerHandler{validate: validate}
}

// SetLevel изменяет уровень логирования, переданный в теле запроса в формате JSON.
// Пример запроса:
//
//	POST /debug/loglevel
//	{ "level": "debug" }
//
// В случае успешного изменения возвращает текущий уровень логирования.
func (h *LoggerHandler) SetLevel(res http.ResponseWriter, req *http.Request) {
	var level model.LogLevel

	if err := json.NewDecoder(req.Body).Decode(&level); err != nil {
		http.Error(res, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(level); err != nil {
		http.Error(res, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logger.SetLevel(level.Level); err != nil {
		logger.Log.Error("failed to set log level", zap.String("level", level.Level), zap.Error(err))
		http.Error(res, "internal error", http.StatusInternalServerError)
		return
	}

	body := []byte(logger.GetLevel())
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(body); err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}

// GetLevel возвращает текущий уровень логирования.
//
// GET /debug/loglevel
func (h *LoggerHandler) GetLevel(res http.ResponseWriter, req *http.Request) {
	body := []byte(logger.GetLevel())
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(body); err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}
