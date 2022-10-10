package service

import (
	"backend/pkg/repository/postgres"
	_ "backend/pkg/repository/postgres"
	"backend/structs"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	CreateUser(ctx context.Context, user structs.User) (int, error)
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
