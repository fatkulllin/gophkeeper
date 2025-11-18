package models

import "net/http"

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Cookies    []*http.Cookie
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// type UserWithKey struct {
// 	User
// }
