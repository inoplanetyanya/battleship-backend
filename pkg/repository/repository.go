package repository

import (
	"battleship/pkg/common"
	"database/sql"
)

type Authorization interface {
	CreateUser(user common.User) (int, error)
	GetUser(username, password string) (common.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{}
}
