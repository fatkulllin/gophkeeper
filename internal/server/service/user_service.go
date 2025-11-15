package service

import (
	"context"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
)

type UserService struct {
	repo         UserRepositories
	password     Password
	tokenManager TokenManager
	cryptoUtil   CryptoUtil
}

func NewUserService(repo UserRepositories, tokenManager TokenManager, password Password, cryptoUtil CryptoUtil) *UserService {
	return &UserService{repo: repo, tokenManager: tokenManager, password: password, cryptoUtil: cryptoUtil}
}

func (s *UserService) UserRegister(ctx context.Context, user model.UserCredentials) (string, int, error) {

	userExists, err := s.repo.ExistUser(ctx, user)

	if err != nil {
		return "", 0, err
	}

	if userExists {
		return "", 0, model.ErrUserExists
	}

	hashPassword, err := s.password.Hash(user.Password)
	if err != nil {
		return "", 0, fmt.Errorf("hash password: %w", err)
	}
	user.Password = hashPassword

	random, err := s.cryptoUtil.GenerateRandom(32)

	if err != nil {
		return "", 0, err
	}

	user.EncryptedKey, err = s.cryptoUtil.EncryptWithMasterKey(random)

	if err != nil {
		return "", 0, err
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", 0, err
	}

	tokenString, tokenExpires, err := s.tokenManager.Generate(userID, user.Username)

	if err != nil {
		return "", 0, err
	}

	return tokenString, tokenExpires, nil
}

func (s *UserService) UserLogin(ctx context.Context, user model.UserCredentials) (string, int, error) {
	getUser, err := s.repo.GetUser(ctx, user)

	if err != nil {
		return "", 0, err
	}
	resultPassword, err := s.password.Compare(getUser.PasswordHash, user.Password)

	if err != nil {
		return "", 0, err
	}

	if !resultPassword {
		return "", 0, model.ErrIncorrectPassword
	}

	tokenString, tokenExpires, err := s.tokenManager.Generate(getUser.ID, getUser.Login)

	if err != nil {
		return "", 0, err
	}
	return tokenString, tokenExpires, nil
}
