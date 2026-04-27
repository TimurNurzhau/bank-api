package repositories

import (
	"database/sql"
	"errors"

	"bank-api/models"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *models.Account) error {
	query := `
		INSERT INTO accounts (user_id, balance, currency)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.db.QueryRow(
		query,
		account.UserID,
		account.Balance,
		account.Currency,
	).Scan(&account.ID, &account.CreatedAt)
}

func (r *AccountRepository) FindByUserID(userID int) ([]models.Account, error) {
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.UserID, &a.Balance, &a.Currency, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

// ИСПРАВЛЕНО: проверка прав прямо в SQL
func (r *AccountRepository) FindByIDAndUserID(id, userID int) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRow(query, id, userID).Scan(
		&account.ID, &account.UserID, &account.Balance,
		&account.Currency, &account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found or access denied")
		}
		return nil, err
	}
	return account, nil
}

// Старый метод оставляем для внутреннего использования (без проверки прав)
func (r *AccountRepository) FindByID(id int) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.UserID, &account.Balance,
		&account.Currency, &account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

func (r *AccountRepository) UpdateBalance(accountID int, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2 AND balance + $1 >= 0`

	result, err := r.db.Exec(query, amount, accountID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("insufficient funds or account not found")
	}
	return nil
}