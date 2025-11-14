package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции Initialize, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()
var atomicLevel zap.AtomicLevel = zap.NewAtomicLevel()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string, developLog bool) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	level = strings.ToLower(level)
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	var cfg zap.Config

	if developLog {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// устанавливаем уровень
	cfg.Level = atomicLevel
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

func SetLevel(level string) error {
	level = strings.ToLower(level)
	return atomicLevel.UnmarshalText([]byte(level))
}

func GetLevel() string {
	return atomicLevel.String()
}
