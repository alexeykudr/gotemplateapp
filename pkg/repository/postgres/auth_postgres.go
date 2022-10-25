package postgres

import (
	"backend/domain"
	"backend/pkg/utils"
	"context"
	"errors"
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

func (i *Instance) GetUser(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	row := i.Db.QueryRow(ctx, "SELECT id, email FROM users WHERE email=$1 AND password_hash=$2;",
		email, password)
	err := row.Scan(&user.ID, &user.Email)

	if err != nil {
		log.Error("Error with scan in GetUser " + err.Error())
		return user, err
	}

	return user, nil
}

func (i *Instance) AddUser(ctx context.Context, user domain.User) (int, error) {
	var id int
	err := i.Db.QueryRow(ctx, "INSERT INTO users (created_at, email, password_hash) VALUES ($1, $2, $3) RETURNING ID",
		time.Now(), user.Email, user.Password).Scan(&id)

	if err != nil {
		log.Error("Error with create user in AddUser" + err.Error())
		return 0, err
	}
	return id, nil
}

func (i *Instance) UpdateUserPassword(ctx context.Context, email string) error {
	newPass := utils.RandStringRunes(8)
	hashedPass := utils.GeneratePasswordHash(newPass)

	tag, err := i.Db.Exec(ctx, "UPDATE users SET password_hash=$1 WHERE email=$2;", hashedPass, email)
	if tag.RowsAffected() == 0 {
		return errors.New("no such user")
	}
	if err == pgx.ErrNoRows {
		log.Error("No such email, error in UpdateUserPassword" + err.Error())
		return err
	}
	if err != nil {
		log.Error("Error with UpdateUserPassword" + err.Error())
		return err
	}
	return nil
}
