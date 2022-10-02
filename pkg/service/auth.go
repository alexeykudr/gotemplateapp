package service

import (
	"backend"
	"backend/pkg/repository/postgres"
	"backend/pkg/utils"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AuthService struct {
	repo *postgres.Repository
}

const tokenTTL = 12 * time.Hour
const signingKey = "kjqwhekzhjk123123bjz"

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repo *postgres.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user backend.User) (int, error) {
	user.Password = utils.GeneratePasswordHash(user.Password)
	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(context.Background(), username, utils.GeneratePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
