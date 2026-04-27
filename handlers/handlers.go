package handlers

import (
	"bank-api/repositories"
	"bank-api/services"

	"github.com/sirupsen/logrus"
)

type Handlers struct {
	Auth      *AuthHandler
	Account   *AccountHandler
	Transfer  *TransferHandler
	Card      *CardHandler
	Credit    *CreditHandler
	Analytics *AnalyticsHandler
}

func NewHandlers(svcs *services.Services, repos *repositories.Repositories, logger *logrus.Logger) *Handlers {
	analyticsService := services.NewAnalyticsService(repos)
	return &Handlers{
		Auth:      NewAuthHandler(svcs.Auth, logger),
		Account:   NewAccountHandler(svcs.Account, logger),
		Transfer:  NewTransferHandler(svcs.Transfer, logger),
		Card:      NewCardHandler(svcs.Card, logger),
		Credit:    NewCreditHandler(svcs.Credit, svcs.CBR, logger),
		Analytics: NewAnalyticsHandler(analyticsService, logger), // <--- ИСПРАВЛЕНО
	}
}