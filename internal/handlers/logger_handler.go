package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/logger"
	"github.com/fatkulllin/gophkeeper/internal/model"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type LoggerHandler struct {
	validdate *validator.Validate
}

func NewLoggerHandler() *LoggerHandler {
	return &LoggerHandler{validdate: validator.New()}
}

func (h *LoggerHandler) SetLevel(res http.ResponseWriter, req *http.Request) {
	var level model.LogLevel

	if err := json.NewDecoder(req.Body).Decode(&level); err != nil {
		http.Error(res, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validdate.Struct(level); err != nil {
		http.Error(res, "validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := logger.SetLevel(level.Level); err != nil {
		logger.Log.Error("zap not set", zap.String("level", level.Level), zap.Error(err))
		http.Error(res, "", http.StatusInternalServerError)
		return
	}

	body := []byte(logger.GetLevel())
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)
	_, err := res.Write(body)
	if err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}

func (h *LoggerHandler) GetLevel(res http.ResponseWriter, req *http.Request) {
	body := []byte(logger.GetLevel())
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)
	_, err := res.Write(body)
	if err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}
