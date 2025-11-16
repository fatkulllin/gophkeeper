package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/server/ctxkeys"
	"github.com/fatkulllin/gophkeeper/logger"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// POST /api/records — добавить новую запись;

// GET /api/records — получить список;

// GET /api/records/{id} — получить запись;

// DELETE /api/records/{id} — удалить.
type RecordService interface {
	Create(ctx context.Context, userID int, input model.RecordInput) error
	GetAll(ctx context.Context, userID int) ([]model.Record, error)
	Get(ctx context.Context, userID int, idRecord string) (model.RecordResponse, error)
	Delete(ctx context.Context, userID, recordID int) error
}
type RecordHandler struct {
	service  RecordService
	validate *validator.Validate
}

func NewRecordHandler(service RecordService) *RecordHandler {
	return &RecordHandler{service: service, validate: validator.New()}
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

	// if err := h.validate.Struct(record); err != nil {
	// 	http.Error(res, "Validation failed: "+err.Error(), http.StatusBadRequest)
	// 	return
	// }

	if err := h.service.Create(req.Context(), claims.UserID, record); err != nil {
		http.Error(res, "error", http.StatusInternalServerError)
	}
}

func (h *RecordHandler) ListRecords(res http.ResponseWriter, req *http.Request) {
	claims, ok := req.Context().Value(ctxkeys.UserContextKey).(model.Claims)

	if !ok {
		http.Error(res, "claims not found", http.StatusUnauthorized)
		return
	}
	result, err := h.service.GetAll(req.Context(), claims.UserID)
	if err != nil {
		http.Error(res, "error", http.StatusInternalServerError)
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(result)
	if err != nil {
		logger.Log.Error("json encoder error", zap.Error(err))
		http.Error(res, "error", http.StatusInternalServerError)
	}
}

func (h *RecordHandler) GetRecord(res http.ResponseWriter, req *http.Request) {
	idRecord := chi.URLParam(req, "id")
	claims, ok := req.Context().Value(ctxkeys.UserContextKey).(model.Claims)

	if !ok {
		http.Error(res, "claims not found", http.StatusUnauthorized)
		return
	}
	record, err := h.service.Get(req.Context(), claims.UserID, idRecord)
	if err != nil {
		http.Error(res, "error", http.StatusInternalServerError)
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(record)
	if err != nil {
		logger.Log.Error("json encoder error", zap.Error(err))
		http.Error(res, "error", http.StatusInternalServerError)
	}
}
