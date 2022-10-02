package postgres

import (
	"backend"
	"backend/pkg/utils"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

type Instance struct {
	Db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *Instance {
	return &Instance{Db: db}
}

func (i *Instance) GetUser(ctx context.Context, username, password string) (backend.User, error) {
	var user backend.User

	row := i.Db.QueryRow(ctx, "SELECT id, username, email, stuff FROM users WHERE username=$1 AND password_hash=$2;",
		username, password)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.IsStuff)

	if err != nil {
		log.Error("Error with scan in GetUser " + err.Error())
		return user, err
	}

	return user, nil
}

func (i *Instance) AddUser(ctx context.Context, user backend.User) (int, error) {
	var id int
	err := i.Db.QueryRow(ctx, "INSERT INTO users (created_at, username, password_hash, email, stuff) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		time.Now(), user.Username, user.Password, user.Email, user.IsStuff).Scan(&id)

	if err != nil {
		log.Error("Error with create user in AddUser" + err.Error())
		return 0, err
	}
	return id, nil
}

func (i *Instance) UpdateUserPassword(ctx context.Context, email string) (string, error) {
	newPass := utils.RandStringRunes(8)
	hashedPass := utils.GeneratePasswordHash(newPass)
	var username string

	//_, err := i.Db.Exec(ctx, "UPDATE users SET password_hash=$1 WHERE email=$2;", hashedPass, email)

	err := i.Db.QueryRow(ctx, "UPDATE users SET password_hash=$1 WHERE email=$2 RETURNING username;", hashedPass, email).Scan(&username)
	if err == pgx.ErrNoRows {
		log.Error("No such email, error in UpdateUserPassword" + err.Error())
		return "", err
	}
	if err != nil {
		log.Error("Error with UpdateUserPassword" + err.Error())
		return "", err
	}
	return newPass, nil
}
