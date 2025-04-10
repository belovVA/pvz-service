package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"pvz-service/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
}

type Service interface {
	UserService
}

type Router struct {
	service Service
}

func NewRouter(service Service, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()
	router := &Router{service: service}
	r.Handle("/register", http.HandlerFunc(router.registerHandler))
	return r
}

func (r *Router) registerHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.Register(w, req)
}
