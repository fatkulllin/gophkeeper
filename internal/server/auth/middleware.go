package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fatkulllin/gophkeeper/internal/server/ctxkeys"
	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// AuthMiddleware проверяет JWT-токен из cookie "auth_token".
// При успешной аутентификации помещает данные пользователя (claims) в контекст
// и передает управление следующему обработчику. В случае ошибки возвращает
// статус 401 Unauthorized.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			cookie, err := req.Cookie("auth_token")

			if err != nil {
				http.Error(res, "unauthorized: missing auth token", http.StatusUnauthorized)
				return
			}

			tokenString := cookie.Value

			claims := model.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, &claims,
				func(t *jwt.Token) (any, error) {
					if t.Method != jwt.SigningMethodHS256 {
						return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
					}
					return []byte(secret), nil
				})

			if err != nil {
				logger.Log.Error("JWT validation failed", zap.Error(err))
				http.Error(res, "unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(res, "unauthorized: token is not valid", http.StatusUnauthorized)
				return
			}

			logger.Log.Debug("JWT token validated", zap.String("login", claims.UserLogin))

			// Передаем claims в контекст запроса
			ctx := context.WithValue(req.Context(), ctxkeys.UserContextKey, claims)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
