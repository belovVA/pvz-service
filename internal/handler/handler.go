package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"pvz-service/internal/model"
)

type AuthService interface {
	Registration(ctx context.Context, user model.User) (*model.User, error)
	Authenticate(ctx context.Context, user model.User) (string, error)
	DummyAuth(ctx context.Context, role string) (string, error)
}

type Service interface {
	AuthService
}

type Router struct {
	service Service
}

func NewRouter(service Service, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()
	router := &Router{service: service}
	r.Post("/register", http.HandlerFunc(router.registerHandler))
	r.Post("/login", http.HandlerFunc(router.loginHandler))

	return r
}

func (r *Router) registerHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.Register(w, req)
}

func (r *Router) loginHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.Login(w, req)
}
