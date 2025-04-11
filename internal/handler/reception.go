package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
	"pvz-service/internal/model"
)

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID) (*model.Reception, error)
}
type ReceptionHandlers struct {
	Service ReceptionService
}

func NewReceptionHandler(service ReceptionService) *ReceptionHandlers {
	return &ReceptionHandlers{
		Service: service,
	}
}

func (h *ReceptionHandlers) OpenNewReception(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReceptionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.PvzID)
	if err != nil {
		pkg.WriteError(w, "Invalid Pvz ID", http.StatusBadRequest)
		return
	}

	recep, err := h.Service.CreateReception(r.Context(), id)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("Failed to create Reception: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToReceptionResponseFromReception(recep)

	pkg.SuccessJSON(w, resp, http.StatusCreated)
}
