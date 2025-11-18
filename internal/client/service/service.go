package service

import (
	"net/http"
	"os"

	"github.com/fatkulllin/gophkeeper/internal/client/models"
	"github.com/fatkulllin/gophkeeper/model"
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
	RemoveFile(filename string) error
}

type Repository interface {
	PutUserKey(userKey string) error
	GetUserKey() ([]byte, error)
	SaveRecords(records []model.Record) error
	Clear() error
	All() ([]model.Record, error)
	Get(id int64) (model.Record, error)
}

type CryptoUtil interface {
	Decrypt(encodedCipher string, key []byte) ([]byte, error)
}

func NewService(apiClient ApiClient, fileManager FileManager, boltDB Repository, cryptoUtil CryptoUtil) *Service {
	return &Service{
		User:   NewUserService(apiClient, fileManager, boltDB),
		Record: NewRecordService(apiClient, fileManager, boltDB, cryptoUtil),
	}
}
