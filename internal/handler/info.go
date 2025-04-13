package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/schema"
	"pvz-service/internal/handler/pkg/response"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

const (
	ErrQueryParameters = "invalid query parameters"
	ErrConvertParams   = "invalid converting query parameters"
	FailedGetPvz       = "Failed to get PVZ"
)

type InfoService interface {
	GetInfoPvz(ctx context.Context, query *model.PvzInfoQuery) ([]*model.Pvz, error)
}

type InfoHandlers struct {
	Service InfoService
}

func NewInfoHandler(service InfoService) *InfoHandlers {
	return &InfoHandlers{
		Service: service,
	}
}

func (h *InfoHandlers) GetInfo(w http.ResponseWriter, r *http.Request) {
	var req dto.PvzInfoRequest
	logger := getLogger(r)

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		response.WriteError(w, ErrQueryParameters, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	pvzInfo, err := converter.ToPvzInfoQueryFromPvzInfoResponse(&req)
	if err != nil {
		response.WriteError(w, ErrConvertParams, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrQueryParameters, slog.String(ErrorKey, err.Error()))
		return
	}

	pvzList, err := h.Service.GetInfoPvz(r.Context(), pvzInfo)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s: %s", FailedGetPvz, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), FailedGetPvz, slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToPvzInfoResponseList(pvzList)
	logger.InfoContext(r.Context(), "successful get info about pvz")

	response.SuccessJSON(w, resp, http.StatusOK)
}
