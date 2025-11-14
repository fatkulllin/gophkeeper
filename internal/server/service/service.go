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
}

type RecordRepository interface {
	CreateRecord(ctx context.Context, record model.Record) error
	GetRecords(ctx context.Context, userID int) ([]model.Record, error)
	DeleteRecord(ctx context.Context, recordID, userID int) error
}

type TokenManager interface {
	Generate(userID int, userLogin string) (string, int, error)
}

type Password interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) (bool, error)
}

func NewService(repo Repositories, tokenManager TokenManager, password Password) *Service {
	return &Service{
		User:   NewUserService(repo, tokenManager, password),
		Record: NewRecordService(repo),
	}
}
