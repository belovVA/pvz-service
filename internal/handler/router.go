package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"pvz-service/internal/middleware"
)

const (
	ErrBodyRequest   = "Invalid Request Body"
	ErrRequestFields = "Invalid Request Fields"
	ErrInvalidRole   = "invalid role in Request"
	ErrUUIDParsing   = "invalid ID format"
)

const (
	ModeratorRole = "moderator"
	EmployeeRole  = "employee"
)

const (
	ElectrType  = "электроника"
	ClothesType = "одежда"
	ShoesType   = "обувь"
)

type Service interface {
	AuthService
	PvzService
	ReceptionService
	ProductService
	InfoService
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

		protected.With(middleware.RequireRoles(ModeratorRole)).Post("/pvz", http.HandlerFunc(router.newPvz))

		protected.With(middleware.RequireRoles(ModeratorRole, EmployeeRole)).Get("/pvz", http.HandlerFunc(router.getInfoPvzByParameters))

		// Cоздаём вложенную группу для ручек, требующих роль employee
		protected.Group(func(emp chi.Router) {
			emp.Use(middleware.RequireRoles(EmployeeRole))
			emp.Post("/receptions", http.HandlerFunc(router.newReception))
			emp.Post("/products", http.HandlerFunc(router.newProduct))
			emp.Post("/pvz/{pvzId}/close_last_reception", http.HandlerFunc(router.closeReception))
			emp.Post("/pvz/{pvzId}/delete_last_product", http.HandlerFunc(router.deleteLastProduct))
		})
	})
	r.Route("/auth", func(r chi.Router) {})
	return r
}

func getValidator(r *http.Request) *validator.Validate {
	if v, ok := r.Context().Value("validator").(*validator.Validate); ok {
		return v
	}
	return validator.New()
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

func (r *Router) newReception(w http.ResponseWriter, req *http.Request) {
	h := NewReceptionHandler(r.service)
	h.OpenNewReception(w, req)
}

func (r *Router) newProduct(w http.ResponseWriter, req *http.Request) {
	h := NewProductHandler(r.service)
	h.CreateNewProduct(w, req)
}

func (r *Router) closeReception(w http.ResponseWriter, req *http.Request) {
	h := NewReceptionHandler(r.service)
	h.CloseLastReception(w, req)
}

func (r *Router) deleteLastProduct(w http.ResponseWriter, req *http.Request) {
	h := NewProductHandler(r.service)
	h.RemoveLastProduct(w, req)
}

func (r *Router) getInfoPvzByParameters(w http.ResponseWriter, req *http.Request) {
	h := NewInfoHandler(r.service)
	h.GetInfo(w, req)
}
