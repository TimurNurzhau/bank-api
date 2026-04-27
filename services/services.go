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
	
	emailService := NewEmailService(cfg.SMTPHost, smtpPort, cfg.SMTPUser, cfg.SMTPPass)

	// Логируем статус PGP
	if cfg.PGPPublicKey == "" {
		logger.Warn("PGP_PUBLIC_KEY not set - card numbers will be stored insecurely (base64)")
	} else {
		logger.Info("PGP_PUBLIC_KEY loaded - card numbers will be encrypted")
	}

	return &Services{
		Auth:     NewAuthService(repos.User, cfg),
		Account:  accountService,
		Transfer: NewTransferService(repos.Account, repos.Transaction, accountService, emailService, repos.User),
		Card:     NewCardService(repos.Card, repos.Account, cfg.HMACSecret, cfg.PGPPublicKey),
		Credit:   NewCreditService(repos.Credit, repos.Account, accountService),
		CBR:      NewCBRService(),
		Email:    emailService,
	}
}