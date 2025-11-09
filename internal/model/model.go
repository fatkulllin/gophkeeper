package model

import (
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
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var ErrUserExists = errors.New("user already exists")
var ErrIncorrectPassword = errors.New("incorrect password")

type User struct {
	ID           int
	Login        string
	PasswordHash string
}

type RecordType string

const (
	RecordTypeLoginPassword RecordType = "login_password"
	RecordTypeText          RecordType = "text"
	RecordTypeBinary        RecordType = "binary"
	RecordTypeCard          RecordType = "card"
)

type Record struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Type      RecordType `json:"type"`
	Metadata  string     `json:"metadata"`
	Data      []byte     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type RecordInput struct {
	Type     RecordType `json:"type" validate:"required"`
	Metadata string     `json:"metadata"`
	Data     []byte     `json:"data" validate:"required"`
}
