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

const (
	MoscowRU = "Москва"
	SpbRU    = "Санкт-Петербург"
	KazanRU  = "Казань"
)

const (
	ErrInvalidCity = "invalid city"
	ErrCreatePvz   = "failed to create PVZ"
)

type PvzService interface {
	AddNewPvz(ctx context.Context, city string) (*model.Pvz, error)
}

type PVZHandlers struct {
	Service PvzService
}

func NewPvzHandler(service PvzService) *PVZHandlers {
	return &PVZHandlers{
		Service: service,
	}
}

func (h *PVZHandlers) CreateNewPvz(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePvzRequest
	logger := getLogger(r)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, ErrBodyRequest, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrBodyRequest, slog.String(ErrorKey, err.Error()))
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		response.WriteError(w, ErrRequestFields, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	if err := validateCity(req.City); err != nil {
		response.WriteError(w, ErrInvalidCity, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrInvalidCity, slog.String(ErrorKey, err.Error()))
		return
	}

	pvz, err := h.Service.AddNewPvz(r.Context(), req.City)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s: %s", ErrCreatePvz, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrCreatePvz, slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToCreatePvzResponseFromPvz(pvz)
	logger.InfoContext(r.Context(), "successful create pvz", slog.String(PvzIDKey, resp.ID))

	response.SuccessJSON(w, resp, http.StatusCreated)
}

func validateCity(city string) error {
	switch city {
	case MoscowRU, SpbRU, KazanRU:
		return nil
	}

	return fmt.Errorf("invalid city: %s", city)
}
