package postgres

import (
	"backend"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Authorization interface {
	GetUser(ctx context.Context, username, password string) (backend.User, error)
	AddUser(ctx context.Context, user backend.User) (int, error)
	UpdateUserPassword(ctx context.Context, email string) (string, error)
}

type Stuff interface {
	GetAllUsers(ctx context.Context) ([]backend.User, error)
}

type Repository struct {
	Authorization
	Stuff
}

func NewRepository(Db *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(Db),
		Stuff:         NewStuffPostgres(Db),
	}
}
