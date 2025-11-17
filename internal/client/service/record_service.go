package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
	"github.com/fatkulllin/gophkeeper/model"
)

type RecordService struct {
	apiClient   ApiClient
	fileManager FileManager
}

func NewRecordService(apiClient ApiClient, fileManager FileManager) *RecordService {
	return &RecordService{
		apiClient:   apiClient,
		fileManager: fileManager,
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
