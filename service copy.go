// package service

// import (
// 	"context"
// 	"fmt"

// 	"github.com/fatkulllin/gophkeeper/model"
// )

// type Repositories interface {
// 	ExistUser(ctx context.Context, user model.UserCredentials) (bool, error)
// 	CreateUser(ctx context.Context, user model.UserCredentials) (int, error)
// 	GetUser(ctx context.Context, user model.UserCredentials) (model.User, error)
// }

// type TokenManager interface {
// 	Generate(userID int, userLogin string) (string, int, error)
// }

// type Password interface {
// 	Hash(password string) (string, error)
// 	Compare(hash string, password string) (bool, error)
// }

// type Service struct {
// 	repo         Repositories
// 	password     Password
// 	tokenManager TokenManager
// }

// func NewService(repo Repositories, tokenManager TokenManager, password Password) *Service {
// 	return &Service{repo: repo, tokenManager: tokenManager, password: password}
// }

// func (s *Service) UserRegister(ctx context.Context, user model.UserCredentials) (string, int, error) {
// 	userExists, err := s.repo.ExistUser(ctx, user)
// 	if err != nil {
// 		return "", 0, err
// 	}

// 	if userExists {
// 		return "", 0, model.ErrUserExists
// 	}

// 	hashPassword, err := s.password.Hash(user.Password)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("hash password: %w", err)
// 	}
// 	user.Password = hashPassword

// 	userID, err := s.repo.CreateUser(ctx, user)
// 	if err != nil {
// 		return "", 0, err
// 	}
// 	tokenString, tokenExpires, err := s.tokenManager.Generate(userID, user.Username)
// 	if err != nil {
// 		return "", 0, err
// 	}

// 	return tokenString, tokenExpires, nil
// }

// func (s *Service) UserLogin(ctx context.Context, user model.UserCredentials) (string, int, error) {
// 	getUser, err := s.repo.GetUser(ctx, user)

// 	if err != nil {
// 		return "", 0, err
// 	}
// 	resultPassword, err := s.password.Compare(getUser.PasswordHash, user.Password)

// 	if err != nil {
// 		return "", 0, err
// 	}

// 	if !resultPassword {
// 		return "", 0, model.ErrIncorrectPassword
// 	}

// 	tokenString, tokenExpires, err := s.tokenManager.Generate(getUser.ID, getUser.Login)

// 	if err != nil {
// 		return "", 0, err
// 	}

// 	return tokenString, tokenExpires, nil
// }

// func (s *Service) CreateRecord(ctx context.Context, userID int, record model.RecordInput) {

// 	// userID, err := s.repo.CreateRecord(ctx, user)
// 	// if err != nil {
// 	// 	return "", 0, err
// 	// }
// 	rec := model.Record{
// 		UserID:   userID,
// 		Type:     record.Type,
// 		Metadata: record.Metadata,
// 		Data:     record.Data,
// 	}
// 	fmt.Println(rec)

// }

func (s *Service) CreateRecord(ctx context.Context, userID int, record model.RecordInput) {

	// userID, err := s.repo.CreateRecord(ctx, user)
	// if err != nil {
	// 	return "", 0, err
	// }
	rec := model.Record{
		UserID:   userID,
		Type:     record.Type,
		Metadata: record.Metadata,
		Data:     record.Data,
	}
	fmt.Println(rec)

}
