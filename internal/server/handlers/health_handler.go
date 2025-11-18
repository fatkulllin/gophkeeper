package handlers

import (
	"net/http"

	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthHTTP(res http.ResponseWriter, req *http.Request) {
	body := []byte("OK")
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)
	_, err := res.Write(body)
	if err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}
