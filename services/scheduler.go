package services

import (
	"time"

	"bank-api/models"
	"bank-api/repositories"

	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	creditRepo   *repositories.CreditRepository
	accountRepo  *repositories.AccountRepository
	emailService *EmailService
	logger       *logrus.Logger
}

func NewScheduler(repos *repositories.Repositories, emailService *EmailService, logger *logrus.Logger) *Scheduler {
	return &Scheduler{
		creditRepo:   repos.Credit,
		accountRepo:  repos.Account,
		emailService: emailService,
		logger:       logger,
	}
}

func (s *Scheduler) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			s.logger.Info("Scheduler: processing overdue payments...")
			s.processOverduePayments()
		}
	}()
	s.logger.Infof("Scheduler started with interval: %v", interval)
}

func (s *Scheduler) processOverduePayments() {
	overduePayments, err := s.creditRepo.FindOverduePayments()
	if err != nil {
		s.logger.Errorf("Scheduler error: %v", err)
		return
	}
	for _, payment := range overduePayments {
		s.processPayment(payment)
	}
}

func (s *Scheduler) processPayment(payment models.PaymentSchedule) {
	penalty := payment.Amount * 0.1
	_ = s.creditRepo.AddPenalty(payment.ID, penalty)
	s.logger.Infof("Penalty added to payment %d: +%.2f", payment.ID, penalty)
}