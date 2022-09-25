package handler

import (
	"backend"
	"backend/pkg/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
	err = h.service.CreateUser(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "ok")
}
func (h *Handler) getUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUserList(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, _ := json.Marshal(users)
	w.Write(body)
}
func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	//userId = r.URL.
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	user, err := h.service.GetUserById(context.Background(), id)

	if user.Name == "" {
		http.Error(w, "no such user", http.StatusBadRequest)
		return
	}
	if err != nil {
		return
	}
	body, _ := json.Marshal(user)
	w.Write(body)
}

func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "WORK!")
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("error with parsing id")
		return
	}
	h.service.DeleteUserById(context.Background(), userId)
}
func FirstMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "First middleware execute")
		next.ServeHTTP(w, r)
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
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.With(FirstMiddleware).Get("/health", h.Healthcheck)
	r.Post("/user", h.createUser)
	r.With(BasicAuthMiddleware).Get("/user", h.getUsersList)
	r.Get("/user/{id:[0-9]+}", h.getUserById)
	r.Delete("/user/{id:[0-9]+}", h.deleteUser)
	//r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)
	return r
}
