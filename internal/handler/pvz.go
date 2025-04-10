package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
	"pvz-service/internal/model"
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	pvz, err := h.Service.AddNewPvz(r.Context(), req.City)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("Failed to create PVZ: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToCreatePvzResponseFromPvz(pvz)

	pkg.SuccessJSON(w, resp, http.StatusCreated)
}
