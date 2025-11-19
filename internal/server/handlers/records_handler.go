// RecordHandler обрабатывает CRUD-операции с пользовательскими записями.
//
// Поддерживаемые эндпоинты:
//
//   - POST   /api/record        — создание записи
//   - GET    /api/records       — получение списка записей
//   - GET    /api/records/{id}  — получение записи по ID
//   - DELETE /api/records/{id}  — удаление записи
//   - PATCH  /api/records/{id}  — обновление записи
//
// Хендлеры извлекают идентификатор пользователя из JWT (через контекст),
// проводят базовую проверку входных данных и вызывают доменный сервис.
// Бизнес-логика и работа с хранилищем находятся в слое service.
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/server/ctxkeys"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// RecordService определяет интерфейс бизнес-логики для операций над
// пользовательскими записями.
type RecordService interface {
	Create(ctx context.Context, userID int, input model.RecordInput) error
	GetAll(ctx context.Context, userID int) ([]model.Record, error)
	Get(ctx context.Context, userID int, idRecord string) (model.RecordResponse, error)
	Delete(ctx context.Context, userID int, idRecord string) error
	Update(ctx context.Context, userID int, idRecord string, record model.RecordUpdateInput) error
}

// RecordHandler обрабатывает HTTP-запросы, связанные с пользовательскими записями.
// Он преобразует входные данные, достаёт идентификатор пользователя из контекста
// и вызывает соответствующие методы RecordService.
type RecordHandler struct {
	service  RecordService
	validate *validator.Validate
}

// NewRecordHandler создаёт новый RecordHandler.
func NewRecordHandler(service RecordService, validate *validator.Validate) *RecordHandler {
	return &RecordHandler{service: service, validate: validate}
}

// CreateRecord обрабатывает создание записи.
//
// POST /api/record
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

	if err := h.service.Create(req.Context(), claims.UserID, record); err != nil {
		http.Error(res, "error", http.StatusInternalServerError)
		return
	}
}

// ListRecords возвращает список всех записей пользователя.
//
// GET /api/records
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

// GetRecord возвращает запись по ID.
//
// GET /api/records/{id}
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

// Delete удаляет запись по ID.
//
// DELETE /api/records/{id}
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

// Update обновляет запись. Разрешено обновлять только те поля,
// которые явно указаны в JSON (metadata, data).
//
// PATCH /api/records/{id}
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
