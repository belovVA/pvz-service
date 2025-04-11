package converter

import (
	"pvz-service/internal/model"
	modelRepo "pvz-service/internal/repository/pgdb/model"
)

func ToProductFromProductRepo(product *modelRepo.Product) *model.Product {
	return &model.Product{
		ID:          product.ID,
		DateTime:    product.DateTime,
		TypeProduct: product.TypeProduct,
		ReceptionID: product.ReceptionID,
	}
}
