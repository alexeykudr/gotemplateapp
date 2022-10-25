package handler

import (
	"backend/domain"
	service "backend/pkg/service"
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
	type mockBehavior func(s *mock_service.MockAuthorization, user domain.User)

	TestTableSignUp := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{name: "Ok",
			inputBody: `{"email": "123@gmail.com", "password":"qwerty123546"}`,
			inputUser: domain.User{
				Email:    "123@gmail.com",
				Password: "qwerty123546",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user domain.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{name: "Wrong email",
			inputBody: `{"email": "abc123", "password": "qwertyabcd"}`,
			inputUser: domain.User{
				Email:    "abc123",
				Password: "qwertyabcd",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user domain.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Error":"email: must be a valid email address."}`,
		},
		{name: "Wrong password",
			inputBody: `{"email":"123@gmail.com", "password":"abc1"}`,
			inputUser: domain.User{Email: "123@gmail.com", Password: "abc1"},
			mockBehavior: func(r *mock_service.MockAuthorization, user domain.User) {
				r.EXPECT().CreateUser(context.Background(), user).Return(0, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"Error":"Password: the length must be between 8 and 50."}`,
		},
	}

	for _, test := range TestTableSignUp {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_service.NewMockAuthorization(ctl)
			test.mockBehavior(repo, test.inputUser)

			service := &service.Service{Authorization: repo}
			h := Handler{Service: service}

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

	t_sigin := []struct {
		name                 string
		inputBody            string
		inputUser            domain.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{name: "ok token",
			inputBody: `"email": "admin@gmail.com", "password": "qwerty123123123"`,
			inputUser: domain.User{
				Email:    "admin@gmail.com",
				Password: "qwerty123123123",
			},
			mockBehavior: func(r *mock_service.MockAuthorization, user domain.User) {
				//token, _ := utils.GenerateToken(user.ID)
				//r.EXPECT().Login("abc", "qwe").Return("ok", nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"ok":true"}`,
		},
	}
	for _, test := range t_sigin {
		t.Run(test.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			repo := mock_service.NewMockAuthorization(ctl)
			test.mockBehavior(repo, test.inputUser)
			service := &service.Service{Authorization: repo}
			h := Handler{Service: service}

			r := chi.NewRouter()
			r.Post("/sign-in", h.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/sign-in",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
