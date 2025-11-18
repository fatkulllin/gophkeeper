package auth

import (
	"fmt"
	"time"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/golang-jwt/jwt/v5"
)

// JWTManager управляет созданием JWT-токенов.
// Он хранит секрет подписи и время жизни токена (в часах).
type JWTManager struct {
	jwtSecret    string
	tokenExpires int // TTL в часах
}

// NewJWTManager создаёт новый менеджер JWT.
func NewJWTManager(jwtSecret string, tokenExpires int) *JWTManager {
	return &JWTManager{
		jwtSecret:    jwtSecret,
		tokenExpires: tokenExpires,
	}
}

// Generate создаёт JWT-токен для заданного пользователя.
// Возвращает строку токена, срок жизни и ошибку при подписи.
func (m *JWTManager) Generate(userID int, userLogin string) (string, int, error) {
	now := time.Now()
	claims := model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(m.tokenExpires) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		UserID:    userID,
		UserLogin: userLogin,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.jwtSecret))
	if err != nil {
		return "", 0, fmt.Errorf("failed to sign jwt: %w", err)
	}

	return tokenString, m.tokenExpires, nil
}
