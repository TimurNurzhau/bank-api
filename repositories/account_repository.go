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

// TransferWithTransaction - перевод между счетами с записью транзакции в одной транзакции
func (r *AccountRepository) TransferWithTransaction(fromID, toID int, amount float64, tx *models.Transaction) error {
	dbTx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	result, err := dbTx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromID)
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

	_, err = dbTx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO transactions (from_account_id, to_account_id, amount, type, status, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err = dbTx.QueryRow(
		query,
		tx.FromAccountID,
		tx.ToAccountID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.Description,
	).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return err
	}

	return dbTx.Commit()
}

// DepositWithTransaction - пополнение счёта с записью транзакции в одной транзакции
func (r *AccountRepository) DepositWithTransaction(accountID int, amount float64, tx *models.Transaction) error {
	dbTx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer dbTx.Rollback()

	_, err = dbTx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, accountID)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO transactions (to_account_id, amount, type, status, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err = dbTx.QueryRow(
		query,
		tx.ToAccountID,
		tx.Amount,
		tx.Type,
		tx.Status,
		tx.Description,
	).Scan(&tx.ID, &tx.CreatedAt)
	if err != nil {
		return err
	}

	return dbTx.Commit()
}
