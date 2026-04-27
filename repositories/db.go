package repositories

import "database/sql"

type Repositories struct {
	User *UserRepository
	DB   *sql.DB
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		DB:   db,
	}
}