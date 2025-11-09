package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/ctxkeys"
	"github.com/fatkulllin/gophkeeper/internal/logger"
	"github.com/fatkulllin/gophkeeper/internal/model"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// POST /api/records — добавить новую запись;

// GET /api/records — получить список;

// GET /api/records/{id} — получить запись;

// DELETE /api/records/{id} — удалить.
type RecordService interface {
	CreateRecord(ctx context.Context, userID int, user model.RecordInput)
	// UserLogin(ctx context.Context, user model.UserCredentials) (string, int, error)
}
type RecordHandler struct {
	service  RecordService
	validate *validator.Validate
}

func NewRecordHandler(service RecordService) *RecordHandler {
	return &RecordHandler{service: service, validate: validator.New()}
	// return &RecordHandler{validate: validator.New()}

}

func (h *RecordHandler) CreateRecord(res http.ResponseWriter, req *http.Request) {
	var record model.RecordInput
	claims, ok := req.Context().Value(ctxkeys.UserContextKey).(model.Claims)

	if !ok {
		http.Error(res, "claims not found", http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(req.Body).Decode(&record); err != nil {
		logger.Log.Error("error decode json", zap.Error(err))
		http.Error(res, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(record); err != nil {
		http.Error(res, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	h.service.CreateRecord(req.Context(), claims.UserID, record)

}
