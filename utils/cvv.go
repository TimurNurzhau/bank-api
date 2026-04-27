// utils/cvv.go
package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// GenerateCVV генерирует безопасный 3-значный CVV код используя crypto/rand
func GenerateCVV() (string, error) {
	// Используем crypto/rand для криптостойкой генерации
	var n uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &n)
	if err != nil {
		return "", fmt.Errorf("failed to generate CVV: %w", err)
	}

	// Берем последние 3 цифры
	cvv := n % 1000

	// Форматируем как 3-значное число с ведущими нулями
	return fmt.Sprintf("%03d", cvv), nil
}
