package handler

import (
	"go-meechok/modules/product/useCase"
)

type (
	productGrpcHandler struct {
		productUsecase useCase.ProductUsecase
	}
)

func NewProductGrpcHandler(productUsecase useCase.ProductUsecase) *productGrpcHandler {
	return &productGrpcHandler{
		productUsecase: productUsecase,
	}
}
