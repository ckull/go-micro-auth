package handler

import (
	"go-meechok/modules/product/model"
	"go-meechok/modules/product/useCase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	ProductHandler interface {
	}

	productHandler struct {
		productUsecase useCase.ProductUsecase
	}
)

func NewProductHandler(productUsecase useCase.ProductUsecase) ProductHandler {
	return &productHandler{
		productUsecase,
	}
}

func (h *productHandler) CreateProduct(c echo.Context) error {
	var createProductReq model.CreateProduct
	if err := c.Bind(&createProductReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	newProduct := &model.Product{
		Name:        createProductReq.Name,
		Description: createProductReq.Description,
		Price:       createProductReq.Price,
		Categories:  createProductReq.Categories,
		IsActive:    createProductReq.IsActive,
		SellerID:    createProductReq.SellerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := h.productUsecase.CreateProduct(newProduct)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, result)
}
