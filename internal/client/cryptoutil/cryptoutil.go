package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type CryptoUtil struct {
	MasterKey []byte
}

func NewCryptoUtil() *CryptoUtil {
	return &CryptoUtil{}
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

func (c *CryptoUtil) Decrypt(encodedCipher string, key []byte) ([]byte, error) {
	fmt.Println(encodedCipher, "aaaaa")
	cipherData, err := base64.StdEncoding.DecodeString(encodedCipher)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}
	fmt.Println("key len=", len(key), "key=", key)
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
