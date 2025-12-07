package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandom returns cryptographically secure random bytes.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	return b, err
}

// Encrypt encrypts data with AES-GCM and returns raw bytes (nonce + ciphertext).
func Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	nonce, err := GenerateRandom(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	out := gcm.Seal(nonce, nonce, data, nil)
	return out, nil
}

// EncryptBase64 encrypts and returns base64 string.
func EncryptBase64(data, key []byte) (string, error) {
	b, err := Encrypt(data, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// Decrypt decrypts raw AES-GCM data (nonce + ciphertext).
func Decrypt(cipherData, key []byte) ([]byte, error) {
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
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// DecryptBase64 accepts base64 string and decrypts it.
func DecryptBase64(encoded string, key []byte) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode base64: %w", err)
	}
	return Decrypt(b, key)
}
