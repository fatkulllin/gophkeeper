package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

type UserService struct {
	apiClient   ApiClient
	fileManager FileManager
	boltDB      Repository
}

func NewUserService(apiClient ApiClient, fileManager FileManager, boltDB Repository) *UserService {
	return &UserService{
		apiClient:   apiClient,
		fileManager: fileManager,
		boltDB:      boltDB,
	}
}

func (s *UserService) LoginUser(ctx context.Context, username, password, url string) (*models.Response, error) {

	user := models.UserRequest{
		Username: username,
		Password: password,
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

	user := models.UserRequest{
		Username: username,
		Password: password,
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

func (s *UserService) SaveUserKey(userKey string) error {
	fmt.Println(userKey)
	err := s.boltDB.PutUserKey(userKey)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ClearDB() error {
	err := s.boltDB.Clear()
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) ClearToken(filename string) error {
	err := s.fileManager.RemoveFile(filename)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) SaveToken(filename string, body string) error {
	permission, err := strconv.ParseUint("0600", 8, 32)

	if err != nil {
		return err
	}

	return s.fileManager.SaveFile(filename, body, os.FileMode(permission))
}
