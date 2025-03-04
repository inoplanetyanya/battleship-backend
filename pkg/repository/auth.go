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
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username, password_hash string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	row := r.db.QueryRow(query, username, password_hash)

	if err := row.Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}

func (r *AuthPostgres) UserExist(username string) (common.User, error) {
	var user common.User

	query := fmt.Sprintf("SELECT id, username FROM %s WHERE username=$1", usersTable)

	row := r.db.QueryRow(query, username)

	if err := row.Scan(&user.Id, &user.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("user not found")
		}
		return user, fmt.Errorf("row scan error: %w", err)
	}

	return user, nil
}
