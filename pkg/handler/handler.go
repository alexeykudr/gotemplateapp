package handler

import (
	"backend"
	"backend/pkg/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var user backend.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	id, err := h.service.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make(map[int]backend.User)
	resp[id] = user
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	type signInInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var input signInInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		return
	}
	token, err := h.service.GenerateToken(input.Username, input.Password)
	if err != nil {
		return
	}
	resp := make(map[string]string)
	resp["token"] = token
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
}

func (h *Handler) RequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.NewUUID()
		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), "request_id", id)))
	})
}
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		a := "admin"
		b := "123"
		if ok {
			if username == a && password == b {
				next.ServeHTTP(writer, request)
			} else {
				writer.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
				http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			}
		}
	})
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/user", h.createUser)
	r.Get("/sign-in", h.signIn)
	//r.With(BasicAuthMiddleware).Get("/user", h.getUsersList)
	//r.Get("/user/{id:[0-9]+}", h.getUserById)
	//r.Delete("/user/{id:[0-9]+}", h.deleteUser)
	////r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)
	return r
}
