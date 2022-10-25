package service

import (
	"backend/domain"
	"backend/pkg/repository/postgres"
	_ "backend/pkg/repository/postgres"
	"context"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	CreateUser(ctx context.Context, user domain.User) (int, error)
	Login(email, password string) (string, error)
}
type Service struct {
	Authorization
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo),
	}
}
