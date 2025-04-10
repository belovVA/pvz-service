package handler

import (
	"encoding/json"
	"net/http"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
)

type AuthHandlers struct {
	Service UserService
}

func NewAuthHandler(service UserService) *AuthHandlers {
	return &AuthHandlers{
		Service: service,
	}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.Service.CreateUser(r.Context(), *converter.ToUserFromCreateUserRequest(&req))
	if err != nil {
		pkg.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := converter.ToCreateUserResponseFromUser(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
