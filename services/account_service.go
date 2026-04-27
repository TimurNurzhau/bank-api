package services

import (
	"errors"

	"bank-api/models"
	"bank-api/repositories"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
}

func NewAccountService(accountRepo *repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(userID int, req *models.CreateAccountRequest) (*models.Account, error) {
	currency := "RUB"
	if req.Currency != "" {
		currency = req.Currency
	}

	account := &models.Account{
		UserID:   userID,
		Balance:  0,
		Currency: currency,
	}

	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) GetUserAccounts(userID int) ([]models.Account, error) {
	return s.accountRepo.FindByUserID(userID)
}

func (s *AccountService) GetAccountByID(accountID, userID int) (*models.Account, error) {
	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return nil, err
	}

	if account.UserID != userID {
		return nil, errors.New("access denied")
	}

	return account, nil
}

func (s *AccountService) Deposit(userID int, req *models.DepositRequest) error {
	account, err := s.GetAccountByID(req.AccountID, userID)
	if err != nil {
		return err
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	return s.accountRepo.UpdateBalance(account.ID, req.Amount)
}

func (s *AccountService) Withdraw(userID int, accountID int, amount float64) error {
	account, err := s.GetAccountByID(accountID, userID)
	if err != nil {
		return err
	}

	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	return s.accountRepo.UpdateBalance(accountID, -amount)
}