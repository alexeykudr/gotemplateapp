package service

import (
	"backend"
	"backend/pkg/repository/postgres"
	_ "backend/pkg/repository/postgres"
	"context"
)

//buisnes logic here
//in service we discribe buisness logic
type Authorization interface {
	CreateUser(ctx context.Context, user backend.User) error
	GenerateToken(username, password string) (string, error)
	GetUserList(ctx context.Context) ([]backend.User, error)
	GetUserById(ctx context.Context, id int) (backend.User, error)
	DeleteUserById(ctx context.Context, id int) error
}
type Service struct {
	Authorization
}

func NewService(repo *postgres.Repository) *Service {
	return &Service{
		Authorization(NewAuthService(repo)),
	}
}
