package service

import (
	"battleship/pkg/common"
	"battleship/pkg/repository"
)

type Authorization interface {
	CreateUser(user common.User) (int, error)
	GetUser(username, password string) (common.User, error)
	UserExist(username string) (common.User, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	GetUserByToken(token string) (common.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.AuthRepository),
	}
}
