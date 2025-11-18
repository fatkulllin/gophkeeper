// UserService реализует бизнес-логику регистрации и авторизации пользователей.
// Он отвечает за хеширование паролей, генерацию и хранение пользовательских ключей,
// шифрование user-key master-key’ем.
package service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// UserService содержит бизнес-логику регистрации и авторизации пользователей.
type UserService struct {
	repo         UserRepositories
	password     Password
	tokenManager TokenManager
	cryptoUtil   CryptoUtil
}

// NewUserService создаёт новый сервис для работы с пользователями
func NewUserService(repo UserRepositories, tokenManager TokenManager, password Password, cryptoUtil CryptoUtil) *UserService {
	return &UserService{repo: repo, tokenManager: tokenManager, password: password, cryptoUtil: cryptoUtil}
}

// UserRegister выполняет регистрацию нового пользователя.
// Генерируется user-key (32 байта), который шифруется master-key’ем,
// пароль хешируется с использованием scrypt, затем создаётся JWT.
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

// UserLogin выполняет авторизацию.
// При wantUserKey = true дополнительно расшифровывает user-key и возвращает его в base64.
func (s *UserService) UserLogin(ctx context.Context, user model.UserCredentials, wantUserKey bool) (string, int, string, error) {
	var userKeyBase64 string
	getUser, err := s.repo.GetUser(ctx, user)
	if wantUserKey {
		encryptedKey, err := s.repo.GetEncryptedKeyUser(ctx, getUser.ID)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return "", 0, "", err
		}
		decryptUserKey, err := s.cryptoUtil.DecryptWithMasterKey(encryptedKey)
		if err != nil {
			logger.Log.Error("", zap.Error(err))
			return "", 0, "", err
		}
		userKeyBase64 = base64.StdEncoding.EncodeToString(decryptUserKey)
	}
	if err != nil {
		return "", 0, "", err
	}
	resultPassword, err := s.password.Compare(getUser.PasswordHash, user.Password)

	if err != nil {
		return "", 0, "", err
	}

	if !resultPassword {
		return "", 0, "", model.ErrIncorrectPassword
	}

	tokenString, tokenExpires, err := s.tokenManager.Generate(getUser.ID, getUser.Login)

	if err != nil {
		return "", 0, "", err
	}
	return tokenString, tokenExpires, userKeyBase64, nil
}
