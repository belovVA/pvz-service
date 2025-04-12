package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg"
	"pvz-service/internal/model"
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

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	if err := decoder.Decode(&req, r.URL.Query()); err != nil {
		pkg.WriteError(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	pvzInfo, err := converter.ToPvzInfoQueryFromPvzInfoResponse(&req)
	if err != nil {
		pkg.WriteError(w, "invalid converting query parameters", http.StatusBadRequest)
		return
	}

	pvzList, err := h.Service.GetInfoPvz(r.Context(), pvzInfo)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("Failed to get PVZ: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToPvzInfoResponseList(pvzList)

	pkg.SuccessJSON(w, resp, http.StatusOK)
}
