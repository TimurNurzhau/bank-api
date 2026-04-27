package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeHMAC(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMAC(data, expectedMAC string, secret []byte) bool {
	computedMAC := ComputeHMAC(data, secret)
	return hmac.Equal([]byte(computedMAC), []byte(expectedMAC))
}