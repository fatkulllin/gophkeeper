package service

import (
	"net/http"
	"os"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
)

type Service struct {
	User   *UserService
	Record *RecordService
}

type ApiClient interface {
	Do(req *http.Request) (*models.Response, error)
}

type FileManager interface {
	SaveFile(filename string, body string, permission os.FileMode) error
	LoadFile(filename string) (string, error)
}

func NewService(apiClient ApiClient, fileManager FileManager) *Service {
	return &Service{
		User:   NewUserService(apiClient, fileManager),
		Record: NewRecordService(apiClient, fileManager),
	}
}
