package service

import (
	"backend"
	mock_service "backend/pkg/service/mocks"
	"bytes"
	"context"
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
			expectedResponseBody: `{"id":3}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_service.NewMockAuthorization(ctl)
			test.mockBehavior(repo, test.inputUser)

			service := &Service{Authorization: repo}
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
