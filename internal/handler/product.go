package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"pvz-service/internal/converter"
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/handler/pkg/response"
	"pvz-service/internal/model"
)

const (
	ErrProductType      = "Invalid Type Product"
	FailedDeleteProduct = "failed to delete product"
	FailedCreateProduct = "Failed add Product"
)

type ProductService interface {
	AddProduct(ctx context.Context, typeProduct string, pvzID uuid.UUID) (*model.Product, error)
	DeleteProduct(ctx context.Context, pvzID uuid.UUID) error
}

type ProductHandlers struct {
	Service ProductService
}

func NewProductHandler(service ProductService) *ProductHandlers {
	return &ProductHandlers{
		Service: service,
	}
}

func (h *ProductHandlers) CreateNewProduct(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
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

	pvzID, err := converter.ParseUuid(req.PvzID)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrRequestFields, slog.String(ErrorKey, err.Error()))
		return
	}

	if err = validateType(req.TypeProduct); err != nil {
		response.WriteError(w, ErrProductType, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrProductType, slog.String(ErrorKey, err.Error()))
		return
	}

	product, err := h.Service.AddProduct(r.Context(), req.TypeProduct, pvzID)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s: %s", FailedCreateProduct, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), FailedCreateProduct, slog.String(ErrorKey, err.Error()))
		return
	}

	resp := converter.ToProductResponseFromProduct(product)
	logger.InfoContext(r.Context(), "successful add product",
		slog.String(ProductIDKey, product.ID.String()),
		slog.String(ReceptionIDKey, product.ID.String()),
	)

	response.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *ProductHandlers) RemoveLastProduct(w http.ResponseWriter, r *http.Request) {
	logger := getLogger(r)
	pvzIdStr := chi.URLParam(r, "pvzId")

	pvzID, err := converter.ParseUuid(pvzIdStr)
	if err != nil {
		response.WriteError(w, ErrUUIDParsing, http.StatusBadRequest)
		logger.InfoContext(r.Context(), ErrUUIDParsing, slog.String(ErrorKey, err.Error()))
		return
	}

	err = h.Service.DeleteProduct(r.Context(), pvzID)
	if err != nil {
		response.WriteError(w, fmt.Sprintf("%s:  %s", FailedDeleteProduct, err.Error()), http.StatusBadRequest)
		logger.InfoContext(r.Context(), FailedDeleteProduct, slog.String(ErrorKey, err.Error()))
		return
	}

	logger.InfoContext(r.Context(), "successful delete last product", slog.String(PvzIDKey, pvzID.String()))
	response.Success(w, http.StatusOK)
}

func validateType(typeProduct string) error {
	switch typeProduct {
	case ElectrType, ClothesType, ShoesType:
		return nil
	}

	return fmt.Errorf("invalid product type: %s", typeProduct)
}
