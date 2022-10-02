package service

import (
	"backend"
	"backend/pkg/repository/postgres"
	_ "backend/pkg/repository/postgres"
	"context"
)

//buisnes logic here
type Authorization interface {
	CreateUser(ctx context.Context, user backend.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}
type Service struct {
	Authorization
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
	}
}