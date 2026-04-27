package repositories

import "database/sql"

type Repositories struct {
	User    *UserRepository
	Account *AccountRepository
	DB      *sql.DB
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Account: NewAccountRepository(db),
		DB:      db,
	}
}