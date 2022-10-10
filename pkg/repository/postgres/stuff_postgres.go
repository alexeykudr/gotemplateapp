package postgres

import (
	"backend/structs"
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

func (i *StuffPostgres) GetAllUsers(ctx context.Context) ([]structs.User, error) {
	var users []structs.User

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
		user := structs.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsStuff)
		if err != nil {
			log.Error("Error with scan in GetAllUsers" + err.Error())
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
