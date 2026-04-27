package services

import (
	"fmt"
	"time"

	"bank-api/repositories"
)

type AnalyticsService struct {
	accountRepo *repositories.AccountRepository
	creditRepo  *repositories.CreditRepository
	txRepo      *repositories.TransactionRepository
}

func NewAnalyticsService(repos *repositories.Repositories) *AnalyticsService {
	return &AnalyticsService{
		accountRepo: repos.Account,
		creditRepo:  repos.Credit,
		txRepo:      repos.Transaction,
	}
}

func (s *AnalyticsService) GetMonthlyStats(userID int) (map[string]float64, error) {
	accounts, err := s.accountRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	stats := map[string]float64{"income": 0, "expenses": 0}
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	for _, acc := range accounts {
		txs, _ := s.txRepo.FindByAccountID(acc.ID, 1000)
		for _, tx := range txs {
			if tx.CreatedAt.After(monthStart) {
				if tx.Type == "deposit" || (tx.ToAccountID != nil && *tx.ToAccountID == acc.ID) {
					stats["income"] += tx.Amount
				}
				if tx.Type == "transfer" && tx.FromAccountID != nil && *tx.FromAccountID == acc.ID {
					stats["expenses"] += tx.Amount
				}
			}
		}
	}
	return stats, nil
}

func (s *AnalyticsService) GetCreditLoad(userID int) (float64, error) {
	credits, err := s.creditRepo.FindByUserID(userID)
	if err != nil {
		return 0, err
	}
	totalMonthly := 0.0
	for _, c := range credits {
		if c.Status == "active" {
			totalMonthly += c.MonthlyPayment
		}
	}
	return totalMonthly, nil
}

func (s *AnalyticsService) PredictBalance(accountID, userID, days int) (float64, error) {
	if days > 365 {
		days = 365
	}

	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return 0, err
	}
	if account.UserID != userID {
		return 0, fmt.Errorf("access denied")
	}

	credits, _ := s.creditRepo.FindByUserID(userID)
	totalPayments := 0.0
	future := time.Now().AddDate(0, 0, days)

	for _, c := range credits {
		if c.Status != "active" {
			continue
		}
		schedule, _ := s.creditRepo.FindScheduleByCreditID(c.ID)
		for _, p := range schedule {
			if p.DueDate.After(time.Now()) && p.DueDate.Before(future) && !p.Paid {
				totalPayments += p.Amount + p.Penalty
			}
		}
	}

	return account.Balance - totalPayments, nil
}