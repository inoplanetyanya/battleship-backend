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
	jwt.RegisteredClaims
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
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

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":  user.Id,
		"username": user.Username,
		"exp":      time.Now().Add(tokenTTL).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signingKey))
}

// TODO rename(?)
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println("parse token error: ", err)
		return -1, err
	}

	if !token.Valid {
		msg := "token is invalid"
		fmt.Println(msg)
		return -1, errors.New(msg)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		fmt.Println("invalid token claims")
		return -1, errors.New("invalid token claims")
	}

	return claims.UserId, nil
}

func (s *AuthService) GetUserByToken(accessToken string) (common.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println("parse token error: ", err)
		return common.User{}, err
	}

	if !token.Valid {
		msg := "token is invalid"
		fmt.Println(msg)
		return common.User{}, errors.New(msg)
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		fmt.Println("invalid token claims")
		return common.User{}, errors.New("invalid token claims")
	}

	user := common.User{
		Id:       claims.UserId,
		Username: claims.Username,
	}

	return user, nil
}
