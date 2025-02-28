package repository

import (
	"battleship/pkg/common"
	"database/sql"
	"errors"
	"fmt"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user common.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	row := r.db.QueryRow(query, username, password)

	if err := row.Scan(&user.Id, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}
