// Package ctxkeys содержит типы и ключи для хранения значений в context.Context.
//
// Пакет используется для передачи данных пользователя (Claims)
// от middleware авторизации к хендлерам.
package ctxkeys

// contextKey — специализированный тип ключей для context.Context.
// Используется для предотвращения коллизий с ключами из других пакетов.
type contextKey string

// UserContextKey — ключ, под которым в контексте хранится информация
// о текущем авторизованном пользователе (структура model.Claims).
const UserContextKey = contextKey("user")
