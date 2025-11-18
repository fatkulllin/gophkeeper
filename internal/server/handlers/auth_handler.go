package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fatkulllin/gophkeeper/model"
	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthService interface {
	UserRegister(ctx context.Context, user model.UserCredentials) (string, int, error)
	UserLogin(ctx context.Context, user model.UserCredentials, wantUserKey bool) (string, int, string, error)
}

type AuthHandler struct {
	service  AuthService
	validate *validator.Validate
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service, validate: validator.New()}
}

func writeAuthSuccessResponse(res http.ResponseWriter, token string, expires int, userKey string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true, // чтобы JS не мог читать cookie (защита от XSS)
		Secure:   true, // true если HTTPS
		Path:     "/",
		MaxAge:   3600 * expires, // время жизни cookie
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(res, cookie)

	if userKey != "" {
		var userKeyResponse model.UserKeyRespone
		userKeyResponse.UserKey = userKey
		res.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(res).Encode(userKeyResponse)
		if err != nil {
			logger.Log.Error("json encoder error", zap.Error(err))
			http.Error(res, "error", http.StatusInternalServerError)
			return
		}
		return
	}

	body := []byte("OK")
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(body); err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}

func (h *AuthHandler) UserRegister(res http.ResponseWriter, req *http.Request) {
	var user model.UserCredentials

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(res, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(user); err != nil {
		http.Error(res, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	tokenString, tokenExpires, err := h.service.UserRegister(req.Context(), user)
	if err != nil {
		if errors.Is(err, model.ErrUserExists) {
			logger.Log.Warn("attempt to register existing user", zap.String("login", user.Username))
			http.Error(res, err.Error(), http.StatusConflict)
			return
		}
		logger.Log.Error("save user", zap.String("login", user.Username), zap.Error(err))
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	writeAuthSuccessResponse(res, tokenString, tokenExpires, "")
}

func (h *AuthHandler) UserLogin(res http.ResponseWriter, req *http.Request) {
	wantUserKey := req.URL.Query().Get("userkey") == "true"

	var user model.UserCredentials

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		logger.Log.Error("", zap.Error(err))
		http.Error(res, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(user); err != nil {
		logger.Log.Error("", zap.Error(err))
		http.Error(res, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	tokenString, tokenExpires, userKey, err := h.service.UserLogin(req.Context(), user, wantUserKey)
	if err != nil {
		if errors.Is(err, model.ErrIncorrectPassword) {
			logger.Log.Warn("attempt to login incorrect password", zap.String("login", user.Username))
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		logger.Log.Error("login user", zap.String("login", user.Username), zap.Error(err))
		http.Error(res, "internal server error", http.StatusInternalServerError)
		return
	}
	writeAuthSuccessResponse(res, tokenString, tokenExpires, userKey)
}

func (h *AuthHandler) UserLogout(res http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true, // чтобы JS не мог читать cookie (защита от XSS)
		Secure:   true, // true если HTTPS
		Path:     "/",
		MaxAge:   -1, // время жизни cookie
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(res, cookie)

	body := []byte("OK")
	res.Header().Set("Content-Type", http.DetectContentType(body))
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(body); err != nil {
		logger.Log.Error("failed to write response", zap.Error(err))
	}
}
