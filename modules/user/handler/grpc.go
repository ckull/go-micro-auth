package handler

// import (
// 	"context"
// 	"go-meechok/modules/inventory/protobuf"
// 	"go-meechok/modules/inventory/useCase"
// )

// type (
// 	inventoryGrpcHandler struct {
// 		inventoryUsecase useCase.InventoryUsecase
// 		protobuf.UnimplementedInventoryServiceServer
// 	}
// )

// func NewProductGrpcHandler(inventoryUsecase useCase.InventoryUsecase) *inventoryGrpcHandler {
// 	return &inventoryGrpcHandler{
// 		inventoryUsecase: inventoryUsecase,
// 	}
// }
// func (g *itemGrpcHandler) FindItemsInIds(ctx context.Context, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
// 	return g.itemUsecase.FindItemInIds(ctx, req)
// }
