package models

import "time"

type Card struct {
	ID               int       `json:"id"`
	AccountID        int       `json:"account_id"`
	EncryptedNumber  string    `json:"-"`
	EncryptedExpiry  string    `json:"-"`
	CVVHash          string    `json:"-"`
	HMAC             string    `json:"-"`
	MaskedNumber     string    `json:"masked_number"`
	ExpiryMonth      int       `json:"expiry_month"`
	ExpiryYear       int       `json:"expiry_year"`
	CreatedAt        time.Time `json:"created_at"`
}