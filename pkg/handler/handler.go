package handler

import (
	"backend/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	"net/http"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(h.RequestIdMiddleware)

	r.Post("/sign-up", h.SignUp)
	r.Get("/sign-in", h.signIn)
	r.With(h.JWTMiddleware).Get("/", h.Healthcheck)
	//r.With(BasicAuthMiddleware).Get("/user", h.getUsersList)
	//r.Get("/user/{id:[0-9]+}", h.getUserById)
	//r.Delete("/user/{id:[0-9]+}", h.deleteUser)
	////r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)
	return r
}
