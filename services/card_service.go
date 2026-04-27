package services

import (
	"errors"
	"fmt"
	"time"

	"bank-api/models"
	"bank-api/repositories"
	"bank-api/utils"
)

type CardService struct {
	cardRepo    *repositories.CardRepository
	accountRepo *repositories.AccountRepository
	hmacSecret  []byte
	pgpKey      string
}

func NewCardService(
	cardRepo *repositories.CardRepository,
	accountRepo *repositories.AccountRepository,
	hmacSecret string,
	pgpKey string,
) *CardService {
	return &CardService{
		cardRepo:    cardRepo,
		accountRepo: accountRepo,
		hmacSecret:  []byte(hmacSecret),
		pgpKey:      pgpKey,
	}
}

func (s *CardService) IssueCard(userID, accountID int) (*models.Card, error) {
	// Проверка владельца счёта
	account, err := s.accountRepo.FindByID(accountID)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Генерация номера карты
	cardNumber := utils.GenerateCardNumber()
	maskedNumber := utils.MaskCardNumber(cardNumber)

	// Срок действия: +3 года
	now := time.Now()
	expiryMonth := int(now.Month())
	expiryYear := now.Year() + 3

	// Генерация CVV
	cvv := fmt.Sprintf("%03d", time.Now().UnixNano()%1000)
	cvvHash, err := utils.HashCVV(cvv)
	if err != nil {
		return nil, err
	}

	// Шифрование номера и срока
	encryptedNumber := utils.EncryptPGP(cardNumber, s.pgpKey)
	expiryStr := fmt.Sprintf("%02d/%d", expiryMonth, expiryYear)
	encryptedExpiry := utils.EncryptPGP(expiryStr, s.pgpKey)

	// HMAC для проверки целостности
	hmac := utils.ComputeHMAC(cardNumber, s.hmacSecret)

	card := &models.Card{
		AccountID:       accountID,
		EncryptedNumber: encryptedNumber,
		EncryptedExpiry: encryptedExpiry,
		CVVHash:         cvvHash,
		HMAC:            hmac,
		MaskedNumber:    maskedNumber,
		ExpiryMonth:     expiryMonth,
		ExpiryYear:      expiryYear,
	}

	if err := s.cardRepo.Create(card); err != nil {
		return nil, err
	}

	return card, nil
}

func (s *CardService) GetUserCards(userID int) ([]models.Card, error) {
	accounts, err := s.accountRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	var allCards []models.Card
	for _, account := range accounts {
		cards, err := s.cardRepo.FindByAccountID(account.ID)
		if err != nil {
			return nil, err
		}
		allCards = append(allCards, cards...)
	}

	return allCards, nil
}

func (s *CardService) GetCardDetails(cardID, userID int) (*models.Card, error) {
	card, err := s.cardRepo.FindByID(cardID)
	if err != nil {
		return nil, err
	}

	// Проверка владельца
	account, err := s.accountRepo.FindByID(card.AccountID)
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, errors.New("access denied")
	}

	// Проверка целостности через HMAC
	// (в реальном проекте — расшифровка и сверка)

	return card, nil
}