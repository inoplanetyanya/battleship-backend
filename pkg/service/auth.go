package service

import (
	"battleship/pkg/common"
	"battleship/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.Claims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user common.User) (int, error) {
	if s.repo == nil {
		return 0, errors.New("repository is not initialized")
	}
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUser(username, password string) (common.User, error) {
	var user common.User

	if s.repo == nil {
		return user, errors.New("repository is not initialized")
	}

	password_hash := generatePasswordHash(password)
	user, err := s.repo.GetUser(username, password_hash)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *AuthService) UserExist(username string) (common.User, error) {
	if s.repo == nil {
		return common.User{}, errors.New("repository is not initialized")
	}

	user, err := s.repo.UserExist(username)
	if err != nil && err.Error() != "user not found" {
		return common.User{}, err
	}

	return user, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
