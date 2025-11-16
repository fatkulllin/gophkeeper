package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

// PATCH /api/records/{id} — обновить.
type RecordService interface {
	Create(ctx context.Context, userID int, input model.RecordInput) error
	GetAll(ctx context.Context, userID int) ([]model.Record, error)
	Get(ctx context.Context, userID int, idRecord string) (model.RecordResponse, error)
	Delete(ctx context.Context, userID int, idRecord string) error
	Update(ctx context.Context, userID int, idRecord string, record model.RecordUpdateInput) error
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
		return
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
		return
	}

	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(result)
	if err != nil {
		logger.Log.Error("json encoder error", zap.Error(err))
		http.Error(res, "error", http.StatusInternalServerError)
		return
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
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(res).Encode(record)
	if err != nil {
		logger.Log.Error("json encoder error", zap.Error(err))
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}
}

func (h *RecordHandler) Delete(res http.ResponseWriter, req *http.Request) {
	idRecord := chi.URLParam(req, "id")
	claims, ok := req.Context().Value(ctxkeys.UserContextKey).(model.Claims)

	if !ok {
		http.Error(res, "claims not found", http.StatusUnauthorized)
		return
	}
	err := h.service.Delete(req.Context(), claims.UserID, idRecord)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(res, "record not found", http.StatusNotFound)
			return
		}
		logger.Log.Error("delete record", zap.String("record id", idRecord), zap.Error(err))
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"status":  "ok",
		"deleted": idRecord,
	})
}

func (h *RecordHandler) Update(res http.ResponseWriter, req *http.Request) {
	var record model.RecordUpdateInput
	idRecord := chi.URLParam(req, "id")
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

	if record.Data == nil && record.Metadata == nil {
		http.Error(res, "no fields to update", http.StatusBadRequest)
		return
	}

	err := h.service.Update(req.Context(), claims.UserID, idRecord, record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(res, "record not found", http.StatusNotFound)
			return
		}
		logger.Log.Error("failed to update record", zap.String("record id", idRecord), zap.Error(err))
		http.Error(res, "failed to update record", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]string{
		"status":  "ok",
		"updated": idRecord,
	})
}
