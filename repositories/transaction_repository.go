package repositories

import (
	"bank-api/models"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *models.Transaction) error {
	query := `
		INSERT INTO transactions (from_account_id, to_account_id, amount, type, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		tx.FromAccountID,
		tx.ToAccountID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.Description,
	).Scan(&tx.ID, &tx.CreatedAt)
}

func (r *TransactionRepository) FindByAccountID(accountID int, limit int) ([]models.Transaction, error) {
	query := `
		SELECT id, from_account_id, to_account_id, amount, type, status, description, created_at
		FROM transactions
		WHERE from_account_id = $1 OR to_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2`

	rows, err := r.db.Query(query, accountID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount,
			&t.Type, &t.Status, &t.Description, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, rows.Err()
}
