package service

import (
	"context"

	"github.com/fatkulllin/gophkeeper/model"
)

type Service struct {
	User   *UserService
	Record *RecordService
}

type Repositories interface {
	UserRepositories
	RecordRepository
}

type UserRepositories interface {
	ExistUser(ctx context.Context, user model.UserCredentials) (bool, error)
	CreateUser(ctx context.Context, user model.UserCredentials) (int, error)
	GetUser(ctx context.Context, user model.UserCredentials) (model.User, error)
	GetEncryptedKeyUser(ctx context.Context, userID int) (string, error)
}

type RecordRepository interface {
	CreateRecord(ctx context.Context, record model.Record) error
	DeleteRecord(ctx context.Context, userID int, idRecord string) error
	GetEncryptedKeyUser(ctx context.Context, userID int) (string, error)
	GetAllRecords(ctx context.Context, userID int) ([]model.Record, error)
	GetRecord(ctx context.Context, userID int, idRecord string) (model.Record, error)
	UpdateRecord(ctx context.Context, userID int, idRecord string, record model.Record) error
}

type TokenManager interface {
	Generate(userID int, userLogin string) (string, int, error)
}

type Password interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) (bool, error)
}

type CryptoUtil interface {
	EncryptWithMasterKey(src []byte) (string, error)
	GenerateRandom(size int) ([]byte, error)
	DecryptWithMasterKey(src string) ([]byte, error)
	EncryptString(src, key []byte) (string, error)
	Decrypt(encodedCipher string, key []byte) ([]byte, error)
}

func NewService(repo Repositories, tokenManager TokenManager, password Password, cryptoUtil CryptoUtil) *Service {
	return &Service{
		User:   NewUserService(repo, tokenManager, password, cryptoUtil),
		Record: NewRecordService(repo, cryptoUtil),
	}
}
