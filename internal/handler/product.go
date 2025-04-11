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

type ProductService interface {
	AddProduct(ctx context.Context, typeProduct string, pvzID uuid.UUID) (*model.Product, error)
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

	pvzID, err := uuid.Parse(req.PvzID)
	if err != nil {
		pkg.WriteError(w, "Invalid Pvz ID", http.StatusBadRequest)
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
