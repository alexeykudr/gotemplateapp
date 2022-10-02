package postgres

import (
	"backend"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type StuffPostgres struct {
	Db *pgxpool.Pool
}

func NewStuffPostgres(db *pgxpool.Pool) *StuffPostgres {
	return &StuffPostgres{Db: db}
}

func (i *StuffPostgres) GetAllUsers(ctx context.Context) ([]backend.User, error) {
	var users []backend.User

	rows, err := i.Db.Query(ctx, "SELECT id, username, password_hash, email, stuff FROM users;")

	if err == pgx.ErrNoRows {
		log.Error("No rows in GetAllUsers" + err.Error())
		return nil, err
	} else if err != nil {
		log.Error("Error in GetAllUsers" + err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		user := backend.User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.IsStuff)
		if err != nil {
			log.Error("Error with scan in GetAllUsers" + err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
