package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

type UserService struct {
	apiClient   ApiClient
	fileManager FileManager
}

func NewUserService(apiClient ApiClient, fileManager FileManager) *UserService {
	return &UserService{
		apiClient:   apiClient,
		fileManager: fileManager,
	}
}

func (s *UserService) LoginUser(ctx context.Context, username, password, url string) (*models.Response, error) {

	user := map[string]string{
		"username": username,
		"password": password,
	}

	reqBody, err := json.Marshal(user)

	if err != nil {
		return nil, fmt.Errorf("failed marshal batch: %w", err)
	}

	bodyReader := bytes.NewBuffer(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}

func (s *UserService) RegisterUser(ctx context.Context, username, password, url string) (*models.Response, error) {
	user := map[string]string{
		"username": username,
		"password": password,
	}

	reqBody, err := json.Marshal(user)

	if err != nil {
		return nil, fmt.Errorf("failed marshal batch: %w", err)
	}

	bodyReader := bytes.NewBuffer(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.apiClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	return resp, nil
}
