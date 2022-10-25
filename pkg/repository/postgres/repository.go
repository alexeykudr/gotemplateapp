package postgres

import (
	"backend/domain"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go
type Authorization interface {
	//поведение авторизоваться, включает getUser AddUser UpdateuserPassword
	GetUser(ctx context.Context, email, password string) (domain.User, error)
	AddUser(ctx context.Context, user domain.User) (int, error)
	UpdateUserPassword(ctx context.Context, email string) error
}

type Stuff interface {
	//задаем поведение став которое включает в себя GetAlluser
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type Repository struct {
	//задаем модель хранения данных Репозитори , в который кладем поведение и можем из нее их как бы брать
	Authorization
	Stuff
}

func NewRepository(Db *pgxpool.Pool) *Repository {
	//Создаем в памяти репозиторий раз эта структура
	return &Repository{
		Authorization: NewAuthPostgres(Db),
	}
}
