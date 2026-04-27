package repositories

import (
	"database/sql"
	"errors"

	"bank-api/models"
)

type CardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *models.Card) error {
	query := `
		INSERT INTO cards (account_id, encrypted_number, encrypted_expiry, cvv_hash, hmac, masked_number, expiry_month, expiry_year)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		card.AccountID,
		card.EncryptedNumber,
		card.EncryptedExpiry,
		card.CVVHash,
		card.HMAC,
		card.MaskedNumber,
		card.ExpiryMonth,
		card.ExpiryYear,
	).Scan(&card.ID, &card.CreatedAt)
}

func (r *CardRepository) FindByAccountID(accountID int) ([]models.Card, error) {
	query := `SELECT id, account_id, encrypted_number, encrypted_expiry, cvv_hash, hmac, masked_number, expiry_month, expiry_year, created_at FROM cards WHERE account_id = $1`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(&c.ID, &c.AccountID, &c.EncryptedNumber, &c.EncryptedExpiry,
			&c.CVVHash, &c.HMAC, &c.MaskedNumber, &c.ExpiryMonth, &c.ExpiryYear, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, rows.Err()
}

func (r *CardRepository) FindByID(id int) (*models.Card, error) {
	card := &models.Card{}
	query := `SELECT id, account_id, encrypted_number, encrypted_expiry, cvv_hash, hmac, masked_number, expiry_month, expiry_year, created_at FROM cards WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&card.ID, &card.AccountID, &card.EncryptedNumber, &card.EncryptedExpiry,
		&card.CVVHash, &card.HMAC, &card.MaskedNumber, &card.ExpiryMonth, &card.ExpiryYear, &card.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("card not found")
		}
		return nil, err
	}
	return card, nil
}