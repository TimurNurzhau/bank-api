package utils

import (
	"testing"
)

func TestPGPEncryption(t *testing.T) {
	// Генерируем ключи
	pub, priv, err := GenerateTestPGPKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	originalText := "4111111111111111"

	// Шифруем
	encrypted := EncryptPGP(originalText, pub)
	if encrypted == "" {
		t.Fatal("Encryption failed")
	}

	// Расшифровываем
	decrypted, err := DecryptPGP(encrypted, priv)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if decrypted != originalText {
		t.Errorf("Expected %s, got %s", originalText, decrypted)
	}

	t.Logf("Original: %s", originalText)
	t.Logf("Encrypted: %s", encrypted[:100]+"...")
	t.Logf("Decrypted: %s", decrypted)
}
