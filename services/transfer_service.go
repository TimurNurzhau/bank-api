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
	emailService    *EmailService
	userRepo        *repositories.UserRepository
}

func NewTransferService(
	accountRepo *repositories.AccountRepository,
	transactionRepo *repositories.TransactionRepository,
	accountService *AccountService,
	emailService *EmailService,
	userRepo *repositories.UserRepository,
) *TransferService {
	return &TransferService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		accountService:  accountService,
		emailService:    emailService,
		userRepo:        userRepo,
	}
}

func (s *TransferService) Transfer(userID int, req *models.TransferRequest) error {
	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}
	if req.FromAccountID == req.ToAccountID {
		return errors.New("cannot transfer to the same account")
	}

	fromAccount, err := s.accountService.GetAccountByID(req.FromAccountID, userID)
	if err != nil {
		return err
	}
	if fromAccount.Balance < req.Amount {
		return errors.New("insufficient funds")
	}

	toAccount, err := s.accountRepo.FindByID(req.ToAccountID)
	if err != nil {
		return errors.New("recipient account not found")
	}

	tx := &models.Transaction{
		FromAccountID: &req.FromAccountID,
		ToAccountID:   &req.ToAccountID,
		Amount:        req.Amount,
		Type:          "transfer",
		Status:        "completed",
		Description:   req.Description,
	}

	if err := s.accountRepo.TransferWithTransaction(req.FromAccountID, req.ToAccountID, req.Amount, tx); err != nil {
		return err
	}

	if fromUser, err := s.userRepo.FindByID(userID); err == nil && fromUser != nil && s.emailService != nil {
		_ = s.emailService.SendPaymentNotification(fromUser.Email, -req.Amount)
	}

	if toUser, err := s.userRepo.FindByID(toAccount.UserID); err == nil && toUser != nil && s.emailService != nil {
		_ = s.emailService.SendPaymentNotification(toUser.Email, req.Amount)
	}

	return nil
}

func (s *TransferService) Deposit(userID int, req *models.DepositRequest) error {
	_, err := s.accountService.GetAccountByID(req.AccountID, userID)
	if err != nil {
		return err
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	tx := &models.Transaction{
		ToAccountID: &req.AccountID,
		Amount:      req.Amount,
		Type:        "deposit",
		Status:      "completed",
		Description: "пополнение счёта",
	}

	if err := s.accountRepo.DepositWithTransaction(req.AccountID, req.Amount, tx); err != nil {
		return err
	}

	if user, err := s.userRepo.FindByID(userID); err == nil && user != nil && s.emailService != nil {
		_ = s.emailService.SendPaymentNotification(user.Email, req.Amount)
	}

	return nil
}