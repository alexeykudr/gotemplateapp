package postgres

import (
	"backend"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"math/rand"
	"time"
)

type Instance struct {
	Db *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *Instance {
	return &Instance{Db: db}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (i *Instance) Start() {
	//
	//u, _ := i.getAllUsers(context.Background())
	//fmt.Println(u)

	//i.updateUserStatus(context.Background(), "user1@gmail.com", true)
	//i.getAllUsers(context.Background())
	//i.getUserByEmail(context.Background(), "user1@gmail.com")
	//i.addLastNameIfIsStuff(context.Background(), "abc")
	//err := i.deleteUserByEmail(context.Background(), "user2@gmail.com")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//

	//for index, o := range StuffUsers {
	//	fmt.Println(index, o)
	//	lastname := RandStringRunes(7)
	//	err := i.updateUserName(context.Background(), o.Name+" "+lastname, o.Email)
	//	if err != nil {
	//		return
	//	}
	//
	//}

}
func (i *Instance) mockUserData() error {
	for j := 0; j < 10; j++ {
		name := RandStringRunes(5)
		username := RandStringRunes(9)
		email := RandStringRunes(15)
		var stuff bool
		s := rand.Intn(2)
		if s == 0 {
			stuff = true
		} else {
			stuff = false
		}

		_, err := i.AddUser(context.Background(), backend.User{
			Name:     name,
			Age:      rand.Intn(40),
			Username: username,
			Email:    email,
			IsStuff:  stuff,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) GetAllUsers(ctx context.Context) ([]backend.User, error) {
	var users []backend.User

	rows, err := i.Db.Query(ctx, "SELECT name, age, username, email, stuff FROM users;")
	if err == pgx.ErrNoRows {
		fmt.Println("No rows")
		return nil, err
	} else if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := backend.User{}
		err := rows.Scan(&user.Name, &user.Age, &user.Username, &user.Email, &user.IsStuff)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
func (i *Instance) GetUserByEmail(ctx context.Context, email string) (backend.User, error) {
	var user backend.User
	err := i.Db.QueryRow(ctx, "SELECT name, age, username, email, stuff FROM users WHERE email=$1;", email).Scan(&user.Name, &user.Age, &user.Username, &user.Email, &user.IsStuff)
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	return user, nil
}

func (i *Instance) GetStuffUsers(ctx context.Context, perms bool) ([]backend.User, error) {
	var users []backend.User

	rows, err := i.Db.Query(ctx, "SELECT name, age, username, email, stuff FROM users WHERE stuff=$1", perms)
	if err == pgx.ErrNoRows {
		return users, nil
	}
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user backend.User
		err := rows.Scan(&user.Name, &user.Age, &user.Username, &user.Email, &user.IsStuff)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (i *Instance) AddUser(ctx context.Context, user backend.User) (int, error) {
	var id int
	err := i.Db.QueryRow(ctx, "INSERT INTO users (created_at, name, age, username, email, stuff) VALUES ($1, $2, $3, $4, $5, $6) RETURNING ID",
		time.Now(), user.Name, user.Age, user.Username, user.Email, user.IsStuff).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (i *Instance) UpdateUserStatus(ctx context.Context, email string, status bool) error {
	_, err := i.Db.Exec(ctx, "UPDATE users SET stuff=$1 WHERE email=$2;", status, email)
	if err == pgx.ErrNoRows {
		fmt.Println("no such email")
	}
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) DeleteUserByEmail(ctx context.Context, email string) error {
	tag, err := i.Db.Exec(ctx, "DELETE FROM users WHERE email=$1;", email)
	if tag.RowsAffected() == 0 {
		fmt.Println("no such email")
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (i *Instance) GetUserById(ctx context.Context, id int) (backend.User, error) {

	var user backend.User
	err := i.Db.QueryRow(ctx, "SELECT name, age, username, email, stuff FROM users WHERE id=$1;", id).Scan(&user.Name, &user.Age, &user.Username, &user.Email, &user.IsStuff)
	if err != nil {
		fmt.Println(err)
		return user, err
	}
	return user, nil
}

func (i *Instance) DeleteUserById(ctx context.Context, id int) error {
	tag, err := i.Db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if tag.RowsAffected() == 0 {
		fmt.Println("no such user")
	}
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
