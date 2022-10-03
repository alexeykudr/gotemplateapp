package postgres

import (
	"backend"
	mock_postgres "backend/pkg/repository/postgres/mocks"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestInstance_AddUser(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	usr1 := backend.User{
		Username: "user1",
		Password: "abc",
		Email:    "qweqew@gmail.com",
		IsStuff:  false,
	}

	repo := mock_postgres.NewMockAuthorization(ctl)

	ctx := context.Background()

	repo.EXPECT().AddUser(ctx, backend.User{
		Username: "user1",
		Password: "abc",
		Email:    "qweqew@gmail.com",
		IsStuff:  false,
	}).Return(100, nil).Times(1)

	//user, err := repo.GetUser(ctx, "user1", "abc")
	//fmt.Println(usr)
	id, _ := repo.AddUser(ctx, usr1)
	fmt.Println(id)

	//if err != nil {
	//	return
	//}

}
