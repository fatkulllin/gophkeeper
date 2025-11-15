package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/fatkulllin/gophkeeper/logger"
	"go.uber.org/zap"
)

type CryptoUtil struct {
	MasterKey []byte
}

func NewCryptoUtil(masterKey string) *CryptoUtil {
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		logger.Log.Fatal("invalid master key", zap.Error(err))
	}
	if len(masterKey) != 32 {
		logger.Log.Fatal("master key must be 32 bytes for AES-256", zap.Int("current", len(masterKey)))
	}
	return &CryptoUtil{MasterKey: key}
}

func (c *CryptoUtil) GenerateRandom(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *CryptoUtil) EncryptString(src, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
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

func (c *CryptoUtil) EncryptWithMasterKey(src []byte) (string, error) {
	return c.EncryptString(src, c.MasterKey)
}

func (c *CryptoUtil) DecryptWithMasterKey(src string) ([]byte, error) {
	res, err := c.Decrypt(src, c.MasterKey)
	if err != nil {
		return []byte{}, fmt.Errorf("decrypt: %w", err)
	}
	return res, nil
}

func (c *CryptoUtil) Decrypt(encodedCipher string, key []byte) ([]byte, error) {
	cipherData, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	if len(cipherData) < gcm.NonceSize() {
		return nil, fmt.Errorf("cipher too short")
	}

	nonce := cipherData[:gcm.NonceSize()]
	ciphertext := cipherData[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}
