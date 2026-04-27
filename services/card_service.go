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
	_, err := s.accountRepo.FindByIDAndUserID(accountID, userID)
	if err != nil {
		return nil, err
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

// Получение деталей карты с проверкой прав
func (s *CardService) GetCardDetails(cardID, userID int) (*models.Card, error) {
	return s.cardRepo.FindByIDAndUserID(cardID, userID)
}

// Оплата картой
func (s *CardService) PayWithCard(cardID, userID int, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// Проверяем права на карту
	card, err := s.cardRepo.FindByIDAndUserID(cardID, userID)
	if err != nil {
		return err
	}

	// Проверяем HMAC (целостность данных)
	if !utils.VerifyHMAC(card.EncryptedNumber, card.HMAC, s.hmacSecret) {
		return errors.New("card data integrity check failed")
	}

	// Получаем счёт
	account, err := s.accountRepo.FindByID(card.AccountID)
	if err != nil {
		return err
	}

	if account.Balance < amount {
		return errors.New("insufficient funds")
	}

	// Списываем деньги
	if err := s.accountRepo.UpdateBalance(card.AccountID, -amount); err != nil {
		return err
	}

	return nil
}