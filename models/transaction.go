package models

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	FromAccountID   *int      `json:"from_account_id,omitempty"`
	ToAccountID     *int      `json:"to_account_id,omitempty"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"`
	Status          string    `json:"status"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransferRequest struct {
	FromAccountID int     `json:"from_account_id"`
	ToAccountID   int     `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

type DepositRequest struct {
	AccountID int     `json:"account_id"`
	Amount    float64 `json:"amount"`
}