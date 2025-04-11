package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
	"pvz-service/internal/model"
)

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error)
	CloseReception(ctx context.Context, pvzID uuid.UUID) (*model.Reception, error)
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
		pkg.WriteError(w, fmt.Sprintf("failed to create Reception: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToReceptionResponseFromReception(recep)

	pkg.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *ReceptionHandlers) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	pvzIdStr := chi.URLParam(r, "pvzId")

	pvzID, err := uuid.Parse(pvzIdStr)
	if err != nil {
		pkg.WriteError(w, "Invalid Pvz ID", http.StatusBadRequest)
		return
	}

	recep, err := h.Service.CloseReception(r.Context(), pvzID)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("failed to close reception:  %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToReceptionResponseFromReception(recep)

	pkg.SuccessJSON(w, resp, http.StatusOK)
}
