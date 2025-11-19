package service

import (
	"context"

	"github.com/fatkulllin/gophkeeper/model"
)

// Service агрегирует все сервисы доменной логики — работу с пользователями и записями.
type Service struct {
	User   *UserService
	Record *RecordService
}

// UserRepositories определяет методы для работы с пользователями в хранилище.
type UserRepositories interface {
	ExistUser(ctx context.Context, user model.UserCredentials) (bool, error)
	CreateUser(ctx context.Context, user model.UserCredentials) (int, error)
	GetUser(ctx context.Context, user model.UserCredentials) (model.User, error)
	GetEncryptedKeyUser(ctx context.Context, userID int) (string, error)
}

// RecordRepository определяет методы работы с записями пользователя.
type RecordRepositories interface {
	CreateRecord(ctx context.Context, record model.Record) error
	DeleteRecord(ctx context.Context, userID int, idRecord string) error
	GetAllRecords(ctx context.Context, userID int) ([]model.Record, error)
	GetRecord(ctx context.Context, userID int, idRecord string) (model.Record, error)
	UpdateRecord(ctx context.Context, userID int, idRecord string, record model.Record) error
}

// TokenManager предоставляет методы генерации JWT-токенов.
type TokenManager interface {
	Generate(userID int, userLogin string) (string, int, error)
}

// Password предоставляет функции хеширования и проверки паролей.
type Password interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) (bool, error)
}

// CryptoUtil предоставляет операции шифрования и расшифровки данных.
type CryptoUtil interface {
	EncryptWithMasterKey(src []byte) (string, error)
	GenerateRandom(size int) ([]byte, error)
	DecryptWithMasterKey(src string) ([]byte, error)
	EncryptString(src, key []byte) (string, error)
	Decrypt(encodedCipher string, key []byte) ([]byte, error)
}

// NewService создаёт контейнер сервисов и связывает бизнес-логику
// с реализациями репозиториев, менеджером токенов, хешированием паролей и криптографией.
func NewService(userRepo UserRepositories, recordRepo RecordRepositories, tokenManager TokenManager, password Password, cryptoUtil CryptoUtil) *Service {
	return &Service{
		User:   NewUserService(userRepo, tokenManager, password, cryptoUtil),
		Record: NewRecordService(recordRepo, userRepo, cryptoUtil),
	}
}
