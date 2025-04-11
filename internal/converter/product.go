package converter

import (
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToProductResponseFromProduct(product *model.Product) *dto.CreateProductResponse {
	return &dto.CreateProductResponse{
		ID:          product.ID.String(),
		DateTime:    product.DateTime,
		TypeProduct: product.TypeProduct,
		ReceptionID: product.ReceptionID.String(),
	}
}
