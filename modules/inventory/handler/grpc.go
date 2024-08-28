package handler

import (
	"go-meechok/modules/inventory/model"
	"go-meechok/modules/inventory/protobuf"
	"go-meechok/modules/inventory/useCase"
)

type (
	inventoryGrpcHandler struct {
		inventoryUsecase useCase.InventoryUsecase
		protobuf.UnimplementedInventoryServiceServer
	}
)

func NewInventoryGrpcHandler(inventoryUsecase useCase.InventoryUsecase) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{
		inventoryUsecase: inventoryUsecase,
	}
}

func (h *inventoryGrpcHandler) AddInventoryProto(req *model.AddInventoryReq) (*protobuf.InventoryInfo, error) {
	return h.inventoryUsecase.AddInventoryProto(req)
}

func (h *inventoryGrpcHandler) FindOneInventoryByProductIdProto(productId string) (*protobuf.InventoryInfo, error) {
	return h.inventoryUsecase.FindOneInventoryByProductIdProto(productId)
}
