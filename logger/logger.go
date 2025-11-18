package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции Initialize, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

// atomicLevel хранит текущий уровень логирования и предоставляет HTTP-обработчик
// для его изменения во время работы приложения.
var atomicLevel zap.AtomicLevel = zap.NewAtomicLevel()

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string, developLog bool) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	level = strings.ToLower(level)
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
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
		return fmt.Errorf("failed to build logger: %w", err)
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// SetLevel изменяет текущий уровень логирования во время работы приложения.
func SetLevel(level string) error {
	level = strings.ToLower(level)
	if err := atomicLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	return nil
}

// GetLevel возвращает текущий уровень логирования в текстовом виде.
func GetLevel() string {
	return atomicLevel.String()
}
