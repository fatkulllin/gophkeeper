package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID    int
	UserLogin string
}

type LogLevel struct {
	Level string `json:"level" validate:"required"`
}

type UserCredentials struct {
	Username     string `json:"username" validate:"required"`
	Password     string `json:"password" validate:"required"`
	EncryptedKey string `json:"omitempty"`
}

var ErrUserExists = errors.New("user already exists")
var ErrIncorrectPassword = errors.New("incorrect password")

// TODO# возращается ошибка 500 когда пользовате не найден
// var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID           int
	Login        string
	PasswordHash string
	EncryptedKey string
}

type RecordType string

const (
	TypeLoginPassword RecordType = "login_password"
	TypeText          RecordType = "text"
	TypeBinary        RecordType = "binary"
	TypeBankCard      RecordType = "bank_card"
)

type Record struct {
	ID        int64      `json:"id"`
	UserID    int        `json:"user_id"`
	Type      RecordType `json:"type"`
	Metadata  string     `json:"metadata,omitempty"`
	Data      []byte     `json:"data,omitempty"` // зашифрованное содержимое
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type RecordInput struct {
	Type     RecordType      `json:"type"`
	Metadata string          `json:"metadata,omitempty"`
	Data     json.RawMessage `json:"data"` // сырые данные без парсинга
}
