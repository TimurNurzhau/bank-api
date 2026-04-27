package services

import (
	"errors"
	"math"
	"time"

	"bank-api/models"
	"bank-api/repositories"
)

type CreditService struct {
	creditRepo  *repositories.CreditRepository
	accountRepo *repositories.AccountRepository
	accountSvc  *AccountService
}

func NewCreditService(
	creditRepo *repositories.CreditRepository,
	accountRepo *repositories.AccountRepository,
	accountSvc *AccountService,
) *CreditService {
	return &CreditService{
		creditRepo:  creditRepo,
		accountRepo: accountRepo,
		accountSvc:  accountSvc,
	}
}

// Расчёт аннуитетного платежа
func calculateAnnuity(amount float64, rate float64, months int) (float64, float64) {
	monthlyRate := rate / 12 / 100

	if monthlyRate == 0 {
		return amount / float64(months), amount
	}

	payment := amount * monthlyRate * math.Pow(1+monthlyRate, float64(months)) / (math.Pow(1+monthlyRate, float64(months)) - 1)
	totalPayment := payment * float64(months)

	return math.Round(payment*100) / 100, math.Round(totalPayment*100) / 100
}

func (s *CreditService) CreateCredit(userID int, req *models.CreateCreditRequest, keyRate float64) (*models.Credit, error) {
	if req.Amount <= 0 {
		return nil, errors.New("amount must be positive")
	}
	if req.TermMonths <= 0 || req.TermMonths > 360 {
		return nil, errors.New("term must be between 1 and 360 months")
	}

	// Ставка = ключевая ставка ЦБ + маржа банка 5%
	rate := keyRate + 5

	monthlyPayment, totalPayment := calculateAnnuity(req.Amount, rate, req.TermMonths)

	credit := &models.Credit{
		UserID:         userID,
		Amount:         req.Amount,
		Rate:           rate,
		TermMonths:     req.TermMonths,
		MonthlyPayment: monthlyPayment,
		TotalPayment:   totalPayment,
		Status:         "active",
	}

	if err := s.creditRepo.Create(credit); err != nil {
		return nil, err
	}

	// Генерация графика платежей
	if err := s.generatePaymentSchedule(credit); err != nil {
		return nil, err
	}

	// Зачисление суммы кредита на счёт
	accounts, err := s.accountRepo.FindByUserID(userID)
	if err != nil || len(accounts) == 0 {
		return nil, errors.New("no account found for credit disbursement")
	}

	if err := s.accountRepo.UpdateBalance(accounts[0].ID, req.Amount); err != nil {
		return nil, err
	}

	return credit, nil
}

func (s *CreditService) generatePaymentSchedule(credit *models.Credit) error {
	for i := 1; i <= credit.TermMonths; i++ {
		dueDate := time.Now().AddDate(0, i, 0)
		schedule := &models.PaymentSchedule{
			CreditID: credit.ID,
			DueDate:  dueDate,
			Amount:   credit.MonthlyPayment,
			Paid:     false,
			Penalty:  0,
		}
		if err := s.creditRepo.CreatePaymentSchedule(schedule); err != nil {
			return err
		}
	}
	return nil
}

func (s *CreditService) GetUserCredits(userID int) ([]models.Credit, error) {
	return s.creditRepo.FindByUserID(userID)
}

func (s *CreditService) GetCreditSchedule(creditID, userID int) ([]models.PaymentSchedule, error) {
	credit, err := s.creditRepo.FindByID(creditID)
	if err != nil {
		return nil, err
	}
	if credit.UserID != userID {
		return nil, errors.New("access denied")
	}

	return s.creditRepo.FindScheduleByCreditID(creditID)
}