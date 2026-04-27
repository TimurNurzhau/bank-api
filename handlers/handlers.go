package handlers

import "bank-api/services"

type Handlers struct {
	Auth     *AuthHandler
	Account  *AccountHandler
	Transfer *TransferHandler
	Card     *CardHandler
	Credit   *CreditHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth:     NewAuthHandler(services.Auth),
		Account:  NewAccountHandler(services.Account),
		Transfer: NewTransferHandler(services.Transfer),
		Card:     NewCardHandler(services.Card),
		Credit:   NewCreditHandler(services.Credit, services.CBR),
	}
}