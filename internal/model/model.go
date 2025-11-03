package model

type LogLevel struct {
	Level string `json:"level" validate:"required"`
}
