package models

import "time"

type Credit struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Amount          float64   `json:"amount"`
	Rate            float64   `json:"rate"`
	TermMonths      int       `json:"term_months"`
	MonthlyPayment  float64   `json:"monthly_payment"`
	TotalPayment    float64   `json:"total_payment"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateCreditRequest struct {
	Amount     float64 `json:"amount"`
	TermMonths int     `json:"term_months"`
}

type PaymentSchedule struct {
	ID        int       `json:"id"`
	CreditID  int       `json:"credit_id"`
	DueDate   time.Time `json:"due_date"`
	Amount    float64   `json:"amount"`
	Paid      bool      `json:"paid"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	Penalty   float64   `json:"penalty"`
}