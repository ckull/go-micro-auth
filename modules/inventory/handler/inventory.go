package handler

import (
	"go-meechok/modules/inventory/useCase"
)

type (
	InventoryHandler interface {
	}

	inventoryHandler struct {
		inventoryUsecase useCase.InventoryUsecase
	}
)

func NewInventoryHandler(inventoryUsecase useCase.InventoryUsecase) InventoryHandler {
	return &inventoryHandler{
		inventoryUsecase,
	}
}
