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

type Game interface {
	AddPlayerToQueue(player common.Player) error
	RemovePlayerFromQueue(player common.Player) (common.Player, error)
	PlayerInQueue(player common.Player) bool
	CreateGameRoom() *common.GameRoom
	GameRoomList() *common.GameList
}

type Service struct {
	Authorization
	Game
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.AuthRepository),
		Game:          NewGameService(),
	}
}
