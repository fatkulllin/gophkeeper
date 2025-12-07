package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/cryptoutil"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

type RecordService struct {
	apiClient   ApiClient
	fileManager FileManager
	boltDB      Repository
}

func NewRecordService(apiClient ApiClient, fileManager FileManager, boltDB Repository) *RecordService {
	return &RecordService{
		apiClient:   apiClient,
		fileManager: fileManager,
		boltDB:      boltDB,
	}
}

func (s *RecordService) Add(ctx context.Context, input model.RecordInput, url string) (*models.Response, error) {

	token, err := s.fileManager.LoadFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed read token: %w", err)
	}
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed marshal batch: %w", err)
	}
	bodyReader := bytes.NewBuffer(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	cookie := &http.Cookie{
		Name:  "auth_token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)
	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}

func (s *RecordService) Get(ctx context.Context, url string) (*models.Response, error) {

	token, err := s.fileManager.LoadFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed read token: %w", err)
	}

	bodyReader := bytes.NewBuffer([]byte{})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	cookie := &http.Cookie{
		Name:  "auth_token",
		Value: token,
		Path:  "/",
	}
	req.AddCookie(cookie)
	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}

func (s *RecordService) GetLocal(ctx context.Context, id int64) (model.RecordResponse, error) {
	record, err := s.boltDB.Get(id)
	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	userKey, err := s.boltDB.GetUserKey()

	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}
	decryptData, err := cryptoutil.Decrypt(record.Data, userKey)

	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return model.RecordResponse{}, err
	}

	return model.RecordResponse{
		ID:       record.ID,
		Type:     record.Type,
		Metadata: record.Metadata,
		Data:     decryptData,
	}, nil

}

func (s *RecordService) Delete(ctx context.Context, url string) (*models.Response, error) {

	token, err := s.fileManager.LoadFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed read token: %w", err)
	}

	bodyReader := bytes.NewBuffer([]byte{})

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	cookie := &http.Cookie{
		Name:  "auth_token",
		Value: token,
		Path:  "/",
	}

	req.AddCookie(cookie)
	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}

func (s *RecordService) Update(ctx context.Context, url string, input model.RecordUpdateInput) (*models.Response, error) {

	token, err := s.fileManager.LoadFile("token")
	if err != nil {
		return nil, fmt.Errorf("failed read token: %w", err)
	}
	reqBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed marshal batch: %w", err)
	}
	bodyReader := bytes.NewBuffer(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	cookie := &http.Cookie{
		Name:  "auth_token",
		Value: token,
		Path:  "/",
	}

	req.AddCookie(cookie)
	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}

func (s *RecordService) SaveRecords(records []model.Record) error {
	if err := s.boltDB.SaveRecords(records); err != nil {
		return err
	}
	return nil
}

func (s *RecordService) GetAll() ([]model.RecordResponse, error) {

	recordsOutput := make([]model.RecordResponse, 0)

	records, err := s.boltDB.All()

	if err != nil {
		logger.Log.Error("failed to get all record", zap.Error(err))
		return nil, err
	}

	userKey, err := s.boltDB.GetUserKey()

	if err != nil {
		logger.Log.Error("", zap.Error(err))
		return nil, err
	}

	for _, rec := range records {
		decryptData, err := cryptoutil.Decrypt(rec.Data, userKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return nil, err
		}
		record := model.RecordResponse{
			ID:       rec.ID,
			Type:     rec.Type,
			Metadata: rec.Metadata,
			Data:     decryptData,
		}

		recordsOutput = append(recordsOutput, record)
	}

	return recordsOutput, nil
}
