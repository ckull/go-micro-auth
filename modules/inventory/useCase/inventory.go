package useCase

import (
	"go-meechok/modules/inventory/model"
	"go-meechok/modules/inventory/protobuf"
	"go-meechok/modules/inventory/repository"
	"time"
)

type (
	InventoryUsecase interface {
		AddInventory(inventory *model.AddInventoryReq) (*model.Inventory, error)
		AddInventoryProto(req *model.AddInventoryReq) (*protobuf.InventoryInfo, error)
		FindOneInventoryByProductIdProto(productId string) (*protobuf.InventoryInfo, error)
	}

	inventoryUsecase struct {
		inventoryRepository repository.InventoryRepository
	}
)

func NewInventoryUsecase(inventoryRepository repository.InventoryRepository) InventoryUsecase {
	return &inventoryUsecase{
		inventoryRepository,
	}
}

func (u *inventoryUsecase) AddInventory(inventory *model.AddInventoryReq) (*model.Inventory, error) {
	return u.inventoryRepository.AddInventory(inventory)
}

func (u *inventoryUsecase) AddInventoryProto(req *model.AddInventoryReq) (*protobuf.InventoryInfo, error) {
	inventory, err := u.inventoryRepository.AddInventory(req)

	if err != nil {
		return nil, err
	}

	inventoryInfo := &protobuf.InventoryInfo{
		Id:        inventory.ID.Hex(),
		ProductId: inventory.ProductID.Hex(),
		Quantity:  int32(inventory.Quantity),
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt.Format(time.RFC3339),
		UpdatedAt: inventory.UpdatedAt.Format(time.RFC3339),
	}

	return inventoryInfo, nil

}

func (u *inventoryUsecase) FindOneInventoryByProductIdProto(productId string) (*protobuf.InventoryInfo, error) {

	inventory, err := u.inventoryRepository.FindOneInventoryByProductId(productId)

	if err != nil {
		return nil, err
	}

	inventoryInfo := &protobuf.InventoryInfo{
		Id:        inventory.ID.Hex(),
		ProductId: inventory.ProductID.Hex(),
		Quantity:  int32(inventory.Quantity),
		Version:   inventory.Version,
		CreatedAt: inventory.CreatedAt.Format(time.RFC3339),
		UpdatedAt: inventory.UpdatedAt.Format(time.RFC3339),
	}

	return inventoryInfo, nil

}
