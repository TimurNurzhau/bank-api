package repositories

import "database/sql"

type Repositories struct {
	User        *UserRepository
	Account     *AccountRepository
	Card        *CardRepository
	Transaction *TransactionRepository
	Credit      *CreditRepository
	DB          *sql.DB
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(db),
		Account:     NewAccountRepository(db),
		Card:        NewCardRepository(db),
		Transaction: NewTransactionRepository(db),
		Credit:      NewCreditRepository(db),
		DB:          db,
	}
}
