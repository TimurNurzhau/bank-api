// utils/pgp.go
package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

// EncryptPGP шифрует данные с использованием публичного PGP-ключа
func EncryptPGP(plaintext string, publicKeyStr string) string {
	if publicKeyStr == "" {
		// Fallback на base64, если ключ не задан
		return base64.StdEncoding.EncodeToString([]byte(plaintext))
	}

	encrypted, err := helper.EncryptMessageArmored(publicKeyStr, plaintext)
	if err != nil {
		// При ошибке возвращаем base64
		return base64.StdEncoding.EncodeToString([]byte(plaintext))
	}
	return encrypted
}

// DecryptPGP расшифровывает данные с использованием приватного PGP-ключа
func DecryptPGP(encrypted string, privateKeyStr string) (string, error) {
	if privateKeyStr == "" {
		decoded, err := base64.StdEncoding.DecodeString(encrypted)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	return helper.DecryptMessageArmored(privateKeyStr, nil, encrypted)
}

// GenerateTestPGPKeys генерирует тестовую пару ключей
func GenerateTestPGPKeys() (publicKey, privateKey string, err error) {
	// Генерируем ключи с помощью gopenpgp
	rsaKey, err := crypto.GenerateKey("Bank API", "bank@example.com", "rsa", 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate key: %v", err)
	}

	publicKey, err = rsaKey.GetArmoredPublicKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to get public key: %v", err)
	}

	privateKey, err = rsaKey.Armor()
	if err != nil {
		return "", "", fmt.Errorf("failed to get private key: %v", err)
	}

	return publicKey, privateKey, nil
}
