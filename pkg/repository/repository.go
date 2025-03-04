package repository

import (
	"battleship/pkg/common"
	"database/sql"
)

type AuthRepository interface {
	CreateUser(user common.User) (int, error)
	GetUser(username, password string) (common.User, error)
	UserExist(username string) (common.User, error)
}

type Repository struct {
	AuthRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		AuthRepository: NewAuthPostgres(db),
	}
}
