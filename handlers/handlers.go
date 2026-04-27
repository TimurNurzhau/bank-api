package handlers

import (
	"bank-api/repositories"
	"bank-api/services"
)

type Handlers struct {
	Auth      *AuthHandler
	Account   *AccountHandler
	Transfer  *TransferHandler
	Card      *CardHandler
	Credit    *CreditHandler
	Analytics *AnalyticsHandler
}

func NewHandlers(svcs *services.Services, repos *repositories.Repositories) *Handlers {
	analyticsService := services.NewAnalyticsService(repos)
	return &Handlers{
		Auth:      NewAuthHandler(svcs.Auth),
		Account:   NewAccountHandler(svcs.Account),
		Transfer:  NewTransferHandler(svcs.Transfer),
		Card:      NewCardHandler(svcs.Card),
		Credit:    NewCreditHandler(svcs.Credit, svcs.CBR),
		Analytics: NewAnalyticsHandler(analyticsService),
	}
}
