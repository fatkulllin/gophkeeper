package cryptoutil

import (
	"encoding/base64"

	"github.com/fatkulllin/gophkeeper/pkg/cryptoutil"
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

// EncryptWithMasterKey шифрует данные с использованием мастер-ключа.
func (c *CryptoUtil) EncryptWithMasterKey(src []byte) (string, error) {
	return cryptoutil.EncryptBase64(src, c.MasterKey)
}

// DecryptWithMasterKey расшифровывает base64-строку с использованием мастер-ключа.
func (c *CryptoUtil) DecryptWithMasterKey(src string) ([]byte, error) {
	return cryptoutil.DecryptBase64(src, c.MasterKey)
}
