package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"pvz-service/internal/middleware"
)

type Service interface {
	AuthService
	PvzService
}

type Router struct {
	service Service
}

func NewRouter(service Service, jwtSecret string) *chi.Mux {
	r := chi.NewRouter()
	router := &Router{service: service}
	r.Use(middleware.NewValidator().Middleware)
	r.Post("/register", http.HandlerFunc(router.registerHandler))
	r.Post("/login", http.HandlerFunc(router.loginHandler))
	r.Post("/dummyLogin", http.HandlerFunc(router.dummyLoginHandler))
	r.Group(func(protected chi.Router) {
		protected.Use(middleware.NewJWT(jwtSecret).Authenticate)
		protected.With(middleware.RequireRoles("moderator")).Post("/pvz", http.HandlerFunc(router.newPvz))
	})
	r.Route("/auth", func(r chi.Router) {})
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

func (r *Router) newPvz(w http.ResponseWriter, req *http.Request) {
	h := NewPvzHandler(r.service)
	h.CreateNewPvz(w, req)
}

func getValidator(r *http.Request) *validator.Validate {
	if v, ok := r.Context().Value("validator").(*validator.Validate); ok {
		return v
	}
	return validator.New()
}
