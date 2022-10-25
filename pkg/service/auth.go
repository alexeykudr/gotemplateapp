package service

import (
	"backend/domain"
	"backend/pkg/repository/postgres"
	"backend/pkg/utils"
	"context"
)

type AuthService struct {
	repo *postgres.Repository
}

func NewAuthService(repo *postgres.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(ctx context.Context, user domain.User) (int, error) {
	user.Password = utils.GeneratePasswordHash(user.Password)

	id, err := s.repo.AddUser(ctx, user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUser(context.Background(), email, utils.GeneratePasswordHash(password))
	if err != nil {
		return "error with getting from repo", err
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "error with get token", err
	}
	return token, nil
}
