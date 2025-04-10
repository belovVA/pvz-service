package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"pvz-service/internal/middleware"
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

func NewRouter(service Service, jwtSecret string, v *middleware.Validator) *chi.Mux {
	r := chi.NewRouter()
	router := &Router{service: service}
	r.Use(v.Middleware)
	r.Post("/register", http.HandlerFunc(router.registerHandler))
	r.Post("/login", http.HandlerFunc(router.loginHandler))
	r.Post("/dummyLogin", http.HandlerFunc(router.dummyLoginHandler))

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

func (r *Router) dummyLoginHandler(w http.ResponseWriter, req *http.Request) {
	h := NewAuthHandler(r.service)
	h.DummyLogin(w, req)
}

func getValidator(r *http.Request) *validator.Validate {
	if v, ok := r.Context().Value("validator").(*validator.Validate); ok {
		return v
	}
	return validator.New()
}
