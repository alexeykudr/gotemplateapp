package postgres

import (
	"backend"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Authorization interface {
	GetAllUsers(ctx context.Context) ([]backend.User, error)
	GetUserByEmail(ctx context.Context, email string) (backend.User, error)
	GetStuffUsers(ctx context.Context, perms bool) ([]backend.User, error)
	AddUser(ctx context.Context, user backend.User) (int, error)
	UpdateUserStatus(ctx context.Context, email string, status bool) error
	DeleteUserByEmail(ctx context.Context, email string) error
	GetUserById(ctx context.Context, id int) (backend.User, error)
	DeleteUserById(ctx context.Context, id int) error
}

type Repository struct {
	Authorization
}

func NewRepository(Db *pgxpool.Pool) *Repository {
	return &Repository{Authorization: NewAuthPostgres(Db)}
}
