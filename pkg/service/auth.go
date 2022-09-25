package service

import (
	"backend"
	"backend/pkg/repository/postgres"
	"context"
)

type AuthService struct {
	repo *postgres.Repository
}

func NewAuthService(repo *postgres.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user backend.User) error {
	_, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetUserList(ctx context.Context) ([]backend.User, error) {
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AuthService) GetUserById(ctx context.Context, id int) (backend.User, error) {
	user, err := s.repo.GetUserById(ctx, id)

	if err != nil {
		return user, nil
	}
	return user, nil
}
func (s *AuthService) DeleteUserById(ctx context.Context, id int) error {
	err := s.repo.DeleteUserById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	return "abc", nil
}
