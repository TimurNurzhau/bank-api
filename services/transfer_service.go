package services

import (
	"errors"

	"bank-api/models"
	"bank-api/repositories"
)

type TransferService struct {
	accountRepo     *repositories.AccountRepository
	transactionRepo *repositories.TransactionRepository
	accountService  *AccountService
}

func NewTransferService(
	accountRepo *repositories.AccountRepository,
	transactionRepo *repositories.TransactionRepository,
	accountService *AccountService,
) *TransferService {
	return &TransferService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		accountService:  accountService,
	}
}

func (s *TransferService) Transfer(userID int, req *models.TransferRequest) error {
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.FromAccountID == req.ToAccountID {
		return errors.New("cannot transfer to the same account")
	}

	// Проверка прав на счёт отправителя
	fromAccount, err := s.accountService.GetAccountByID(req.FromAccountID, userID)
	if err != nil {
		return err
	}

	if fromAccount.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	// Проверка существования счёта получателя
	_, err = s.accountRepo.FindByID(req.ToAccountID)
	if err != nil {
		return errors.New("recipient account not found")
	}

	// Атомарная операция: списание + зачисление
	if err := s.accountRepo.UpdateBalance(req.FromAccountID, -req.Amount); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateBalance(req.ToAccountID, req.Amount); err != nil {
		// Откат списания
		_ = s.accountRepo.UpdateBalance(req.FromAccountID, req.Amount)
		return err
	}

	// Запись транзакции
	tx := &models.Transaction{
		FromAccountID: &req.FromAccountID,
		ToAccountID:   &req.ToAccountID,
		Amount:        req.Amount,
		Type:          "transfer",
		Status:        "completed",
		Description:   req.Description,
	}
	_ = s.transactionRepo.Create(tx)

	return nil
}

func (s *TransferService) Deposit(userID int, req *models.DepositRequest) error {
	if err := s.accountService.Deposit(userID, req); err != nil {
		return err
	}

	tx := &models.Transaction{
		ToAccountID: &req.AccountID,
		Amount:      req.Amount,
		Type:        "deposit",
		Status:      "completed",
		Description: "пополнение счёта",
	}
	_ = s.transactionRepo.Create(tx)

	return nil
}