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
	userRepo     *repositories.UserRepository
	emailService *EmailService
	logger       *logrus.Logger
}

func NewScheduler(repos *repositories.Repositories, emailService *EmailService, logger *logrus.Logger) *Scheduler {
	return &Scheduler{
		creditRepo:   repos.Credit,
		accountRepo:  repos.Account,
		userRepo:     repos.User,
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
		s.logger.Errorf("Scheduler error finding overdue payments: %v", err)
		return
	}

	for _, payment := range overduePayments {
		s.processPayment(payment)
	}
}

func (s *Scheduler) processPayment(payment models.PaymentSchedule) {
	// Получаем кредит
	credit, err := s.creditRepo.FindByID(payment.CreditID)
	if err != nil {
		s.logger.Errorf("Scheduler: cannot find credit %d: %v", payment.CreditID, err)
		return
	}

	// Получаем пользователя
	user, err := s.userRepo.FindByID(credit.UserID)
	if err != nil || user == nil {
		s.logger.Errorf("Scheduler: cannot find user for credit %d: %v", credit.ID, err)
		return
	}

	// Получаем счета пользователя
	accounts, err := s.accountRepo.FindByUserID(user.ID)
	if err != nil || len(accounts) == 0 {
		s.logger.Errorf("Scheduler: no accounts found for user %d", user.ID)
		return
	}

	account := accounts[0]
	totalAmount := payment.Amount + payment.Penalty

	// Пытаемся списать
	if account.Balance >= totalAmount {
		// Списание
		if err := s.accountRepo.UpdateBalance(account.ID, -totalAmount); err != nil {
			s.logger.Errorf("Scheduler: failed to withdraw from account %d: %v", account.ID, err)
			return
		}

		// Отмечаем платеж оплаченным
		if err := s.creditRepo.MarkPaymentPaid(payment.ID, time.Now()); err != nil {
			s.logger.Errorf("Scheduler: failed to mark payment %d as paid: %v", payment.ID, err)
		}

		// Отправляем email
		if s.emailService != nil {
			_ = s.emailService.SendPaymentNotification(user.Email, totalAmount)
		}

		s.logger.Infof("Scheduler: successfully processed payment %d for user %d, amount %.2f", payment.ID, user.ID, totalAmount)
	} else {
		// Недостаточно средств - начисляем штраф
		penalty := payment.Amount * 0.1
		if err := s.creditRepo.AddPenalty(payment.ID, penalty); err != nil {
			s.logger.Errorf("Scheduler: failed to add penalty to payment %d: %v", payment.ID, err)
		}

		// Отправляем email о просрочке
		if s.emailService != nil {
			_ = s.emailService.SendCreditReminder(user.Email, payment.Amount, payment.DueDate.Format("2006-01-02"))
		}

		s.logger.Warnf("Scheduler: insufficient funds for payment %d, added penalty %.2f", payment.ID, penalty)
	}
}