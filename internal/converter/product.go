package converter

import (
	"pvz-service/internal/handler/dto"
	"pvz-service/internal/model"
)

func ToProductResponseFromProduct(product *model.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:          product.ID.String(),
		DateTime:    product.DateTime,
		TypeProduct: product.TypeProduct,
		ReceptionID: product.ReceptionID.String(),
	}
}

func ToProductFromCreateProductRequest(request *dto.CreateProductRequest) *model.Product {
	return &model.Product{
		TypeProduct: request.TypeProduct,
	}
}
