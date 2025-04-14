package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg/response"
	"pvz-service/internal/model"
)

type AuthService interface {
	Registration(ctx context.Context, user model.User) (*model.User, error)
	Authenticate(ctx context.Context, user model.User) (string, error)
	DummyAuth(ctx context.Context, user model.User) (string, error)
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
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info(ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info(ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	userModel := *converter.ToUserFromCreateUserRequest(&req)
	if err := validateRole(userModel.Role); err != nil {
		response.WriteError(w, ErrInvalidRole, http.StatusBadRequest)
		logger.Info(ErrInvalidRole, slog.String(ErrorKey, err.Error()))
		return
	}

	user, err := h.Service.Registration(r.Context(), userModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error to register user", slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToCreateUserResponseFromUser(user)
	logger.InfoContext(r.Context(), "successful register", slog.String(UserIDKey, user.ID.String()))

	response.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginUserRequest
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusUnauthorized)
		logger.Info(ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusUnauthorized)
		logger.Info(ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	token, err := h.Service.Authenticate(r.Context(), *converter.ToUserFromLoginUserRequest(&req))
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusUnauthorized)
		logger.Info("error to login user", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.InfoContext(r.Context(), "successful login", slog.String("email", req.Email))

	response.SuccessText(w, token, http.StatusOK)
}

func (h *AuthHandlers) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.TestUserRequest
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.Info(ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.Info(ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	userModel := *converter.ToUserFromDummyLoginRequest(&req)
	if err := validateRole(userModel.Role); err != nil {
		response.WriteError(w, ErrInvalidRole, http.StatusBadRequest)
		logger.Info(ErrInvalidRole, slog.String(ErrorKey, err.Error()))
		return
	}

	token, err := h.Service.DummyAuth(r.Context(), userModel)
	if err != nil {
		response.WriteError(w, err.Error(), http.StatusBadRequest)
		logger.Info("error to login testUser", slog.String(ErrorKey, err.Error()))
		return
	}

	logger.InfoContext(r.Context(), "successful dummyLogin", slog.String("role", req.Role))

	response.SuccessText(w, token, http.StatusOK)
}

func validateRole(role string) error {
	switch role {
	case EmployeeRole, ModeratorRole:
		return nil
	}

	return fmt.Errorf(ErrInvalidRole)
}
