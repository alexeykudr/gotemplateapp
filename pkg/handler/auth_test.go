package handler

import (
	"backend"
	service "backend/pkg/service"
	mock_service "backend/pkg/service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user backend.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            backend.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{name: "Ok",
			inputBody: `{"username": "username", "password": "qwerty", "email": "123@gmail.com", "isStuff": false}`,
			inputUser: backend.User{
				Username: "username",
				Password: "qwerty",
				Email:    "123@gmail.com",
				IsStuff:  false,
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user backend.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{name: "OK",
			inputBody: `{"username": "username", "password": "qwerty", "email": "123@gmail.com"}`,
			inputUser: backend.User{
				Username: "username",
				Password: "qwerty",
				Email:    "123@gmail.com",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user backend.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(2, nil).AnyTimes()
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":2}`,
		},
		{name: "shor username",
			inputBody: `{"username": "q", "password": "qwerty"}`,
			inputUser: backend.User{
				Username: "q",
				Password: "qwerty",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user backend.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, errors.New(`{"message":"Something went wrong"}`)).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Username":"the length must be between 5 and 50"}`,
		},
		{name: "without email",
			inputBody: `{"username":"user123", "password":"abc123456789"}`,
			inputUser: backend.User{
				Username: "user123",
				Password: "abc123456789",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user backend.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, errors.New(`{"message":"Something went wrong"}`)).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"email":"cannot be blank"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_service.NewMockAuthorization(ctl)
			test.mockBehavior(repo, test.inputUser)

			service := &service.Service{Authorization: repo}
			h := Handler{service}

			r := chi.NewRouter()
			r.Post("/sign-up", h.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
