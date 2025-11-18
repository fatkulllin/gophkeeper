package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
)

// CryptoUtil предоставляет функции шифрования и расшифровки данных
// с использованием алгоритма AES-256-GCM.
type CryptoUtil struct {
	MasterKey []byte
}

// NewCryptoUtil принимает base64-представление мастер-ключа,
// декодирует его и проверяет, что длина ключа равна 32 байтам (AES-256).
func NewCryptoUtil(masterKey string) *CryptoUtil {
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		logger.Log.Fatal("invalid master key (must be base64-encoded)", zap.Error(err))
	}
	if len(key) != 32 {
		logger.Log.Fatal("master key must be 32 bytes for AES-256", zap.Int("current", len(key)))
	}
	return &CryptoUtil{MasterKey: key}
}

// GenerateRandom возвращает криптографически случайную последовательность байтов заданного размера.
func (c *CryptoUtil) GenerateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return b, nil
}

// EncryptString шифрует данные с помощью AES-256-GCM и произвольного ключа.
// Возвращает результат в виде base64-строки.
func (c *CryptoUtil) EncryptString(src, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	nonce, err := c.GenerateRandom(gcm.NonceSize())
	if err != nil {
		return "", err
	}

	// шифруем: результат = nonce + ciphertext + tag
	ciphertext := gcm.Seal(nonce, nonce, src, nil)
	// кодируем в base64 для хранения в БД
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptWithMasterKey шифрует данные с использованием мастер-ключа.
func (c *CryptoUtil) EncryptWithMasterKey(src []byte) (string, error) {
	return c.EncryptString(src, c.MasterKey)
}

// DecryptWithMasterKey расшифровывает base64-строку с использованием мастер-ключа.
func (c *CryptoUtil) DecryptWithMasterKey(src string) ([]byte, error) {
	res, err := c.Decrypt(src, c.MasterKey)
	if err != nil {
		return []byte{}, fmt.Errorf("decrypt: %w", err)
	}
	return res, nil
}

// Decrypt расшифровывает base64-строку, зашифрованную AES-256-GCM.
func (c *CryptoUtil) Decrypt(encodedCipher string, key []byte) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(cipherData) < gcm.NonceSize() {
		return nil, fmt.Errorf("cipher data too short")
	}

	nonce := cipherData[:gcm.NonceSize()]
	ciphertext := cipherData[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}
