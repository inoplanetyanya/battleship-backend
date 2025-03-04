package repository

import (
	"database/sql"
	"fmt"
)

const (
	usersTable = "users"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBname   string
	SSLmode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.DBname, cfg.Password, cfg.SSLmode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
