package services

import (
	"log"
	"time"

	"bank-api/models"
	"bank-api/repositories"
)

type Scheduler struct {
	creditRepo   *repositories.CreditRepository
	accountRepo  *repositories.AccountRepository
	emailService *EmailService
}

func NewScheduler(repos *repositories.Repositories, emailService *EmailService) *Scheduler {
	return &Scheduler{
		creditRepo:   repos.Credit,
		accountRepo:  repos.Account,
		emailService: emailService,
	}
}

func (s *Scheduler) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			log.Println("Scheduler: processing overdue payments...")
			s.processOverduePayments()
		}
	}()

	log.Printf("Scheduler started with interval: %v", interval)
}

func (s *Scheduler) processOverduePayments() {
	overduePayments, err := s.creditRepo.FindOverduePayments()
	if err != nil {
		log.Printf("Scheduler error: %v", err)
		return
	}

	for _, payment := range overduePayments {
		s.processPayment(payment)
	}
}

func (s *Scheduler) processPayment(payment models.PaymentSchedule) {
	penalty := payment.Amount * 0.1
	_ = s.creditRepo.AddPenalty(payment.ID, penalty)
	log.Printf("Penalty added to payment %d: +%.2f", payment.ID, penalty)
}