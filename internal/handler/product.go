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

const (
	ElectrType  = "электроника"
	ClothesType = "одежда"
	ShoesType   = "обувь"
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	v := getValidator(r)
	if err := v.Struct(req); err != nil {
		pkg.WriteError(w, "Invalid Request Fields", http.StatusBadRequest)
		return
	}

	pvzID, err := converter.ParseUuid(req.PvzID)
	if err != nil {
		pkg.WriteError(w, "Invalid Pvz ID", http.StatusBadRequest)
		return
	}

	if err = validateType(req.TypeProduct); err != nil {
		pkg.WriteError(w, "Invalid Type Product", http.StatusBadRequest)
		return
	}

	product, err := h.Service.AddProduct(r.Context(), req.TypeProduct, pvzID)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("Failed to create Product: %s", err.Error()), http.StatusBadRequest)
		return
	}

	resp := converter.ToProductResponseFromProduct(product)

	pkg.SuccessJSON(w, resp, http.StatusCreated)
}

func (h *ProductHandlers) RemoveLastProduct(w http.ResponseWriter, r *http.Request) {
	pvzIdStr := chi.URLParam(r, "pvzId")

	pvzID, err := uuid.Parse(pvzIdStr)
	if err != nil {
		pkg.WriteError(w, "Invalid Pvz ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteProduct(r.Context(), pvzID)
	if err != nil {
		pkg.WriteError(w, fmt.Sprintf("failed to delete product:  %s", err.Error()), http.StatusBadRequest)
		return
	}

	pkg.Success(w, http.StatusOK)
}

func validateType(typeProduct string) error {
	switch typeProduct {
	case ElectrType, ClothesType, ShoesType:
		return nil
	}

	return fmt.Errorf("invalid product type: %s", typeProduct)
}
