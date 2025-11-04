package model

import (
	"errors"

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
