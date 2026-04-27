package services

import (
	"strconv"

	"bank-api/config"
	"bank-api/repositories"

	"github.com/sirupsen/logrus"
)

type Services struct {
	Auth     *AuthService
	Account  *AccountService
	Transfer *TransferService
	Card     *CardService
	Credit   *CreditService
	CBR      *CBRService
	Email    *EmailService
}

func NewServices(repos *repositories.Repositories, cfg *config.Config, logger *logrus.Logger) *Services {
	accountService := NewAccountService(repos.Account)
	smtpPort, _ := strconv.Atoi(cfg.SMTPPort)

	return &Services{
		Auth:     NewAuthService(repos.User, cfg),
		Account:  accountService,
		Transfer: NewTransferService(repos.Account, repos.Transaction, accountService),
		Card:     NewCardService(repos.Card, repos.Account, cfg.HMACSecret, cfg.PGPKey),
		Credit:   NewCreditService(repos.Credit, repos.Account, accountService),
		CBR:      NewCBRService(),
		Email:    NewEmailService(cfg.SMTPHost, smtpPort, cfg.SMTPUser, cfg.SMTPPass),
	}
}