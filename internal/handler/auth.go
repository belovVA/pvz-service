package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
	"pvz-service/internal/model"
)

type AuthService interface {
	Registration(ctx context.Context, user model.User) (*model.User, error)
	Authenticate(ctx context.Context, user model.User) (string, error)
	DummyAuth(ctx context.Context, role string) (string, error)
}

type AuthHandlers struct {
	Service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandlers {
	return &AuthHandlers{
		Service: service,
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	user, err := h.Service.Registration(r.Context(), *converter.ToUserFromCreateUserRequest(&req))
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := converter.ToCreateUserResponseFromUser(user)

	pkg.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body: "+err.Error(), http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Authenticate(r.Context(), *converter.ToUserFromLoginUserRequest(&req))
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	pkg.SuccessText(w, token, http.StatusOK)
}

func (h *AuthHandlers) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.TestUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	token, err := h.Service.DummyAuth(r.Context(), req.Role)
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	pkg.SuccessText(w, token, http.StatusOK)
}
