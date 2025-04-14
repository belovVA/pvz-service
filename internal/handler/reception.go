package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg/response"
	"pvz-service/internal/model"
)

const (
	FailedCreateReception = "failed to create Reception"
	FailedCloseReception  = "failed to close reception"
)

type ReceptionService interface {
	CreateReception(ctx context.Context, reception model.Reception) (*model.Reception, error)
	CloseReception(ctx context.Context, reception model.Reception) (*model.Reception, error)
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
	var req dto.ReceptionRequest
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

	receptionModel, err := converter.ToReceptionFromReceptionRequest(&req)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		return
	}

	recep, err := h.Service.CreateReception(r.Context(), *receptionModel)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s: %s", FailedCreateReception, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), FailedCreateReception, slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToReceptionResponseFromReception(recep)
	logger.InfoContext(r.Context(), "successful create reception", slog.String(PvzIDKey, resp.ID))

	response.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *ReceptionHandlers) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	pvzIdStr := chi.URLParam(r, "pvzId")
	logger := getLogger(r)

	receptionModel, err := converter.ToReceptionFromPvzIDRequest(pvzIdStr)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		return
	}

	recep, err := h.Service.CloseReception(r.Context(), *receptionModel)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s: %s", FailedCloseReception, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), FailedCloseReception, slog.String(ErrorKey, err.Error()))

		return
	}

	resp := converter.ToReceptionResponseFromReception(recep)
	logger.InfoContext(r.Context(), "successful close reception", slog.String(PvzIDKey, resp.ID))

	response.SuccessJSON(w, resp, http.StatusOK)
}
