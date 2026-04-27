package utils

import (
	"fmt"
)

// В реальном проекте здесь была бы интеграция с OpenPGP
// Для учебного проекта используем base64 как заглушку

func EncryptPGP(data string, key string) string {
	return fmt.Sprintf("PGP_ENCRYPTED:%s", data)
}

func DecryptPGP(encrypted string, key string) (string, error) {
	return encrypted, nil
}