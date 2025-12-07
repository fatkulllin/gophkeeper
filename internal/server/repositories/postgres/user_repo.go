package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fatkulllin/gophkeeper/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// ExistUser проверяет наличие пользователя по логину.
// Возвращает true, если пользователь существует.
func (s *UserRepo) ExistUser(ctx context.Context, user model.UserCredentials) (bool, error) {
	row := s.db.QueryRowContext(ctx, "SELECT login FROM users WHERE login = $1", user.Username)
	var userScan string
	err := row.Scan(&userScan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check existing user: %w", err)
	}
	return true, nil
}

// CreateUser создаёт нового пользователя и возвращает его ID.
func (s *UserRepo) CreateUser(ctx context.Context, user model.UserCredentials) (int, error) {

	var id int

	row := s.db.QueryRowContext(ctx, "INSERT INTO users (login, password_hash, encrypted_key) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Password, user.EncryptedKey)

	err := row.Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("pg failed to insert new user: %w", err)
	}

	return id, nil
}

// GetUser возвращает пользователя по логину.
func (s *UserRepo) GetUser(ctx context.Context, user model.UserCredentials) (model.User, error) {
	var foundUser model.User
	row := s.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users WHERE login = $1", user.Username)
	err := row.Scan(&foundUser.ID, &foundUser.Login, &foundUser.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found")
		}
		return model.User{}, err
	}

	return foundUser, nil
}

// GetEncryptedKeyUser возвращает зашифрованный ключ пользователя.
func (s *UserRepo) GetEncryptedKeyUser(ctx context.Context, userID int) (string, error) {
	var encryptedKey string
	row := s.db.QueryRowContext(ctx, "SELECT encrypted_key FROM users WHERE id = $1;", userID)
	err := row.Scan(&encryptedKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return encryptedKey, nil
}
