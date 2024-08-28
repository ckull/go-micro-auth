package route

import (
	"go-meechok/modules/inventory/handler"
	"go-meechok/modules/inventory/repository"
	"go-meechok/modules/inventory/useCase"
	"go-meechok/pkg/redisCon"
	"go-meechok/server/types"
)

func InventoryRoute(s *types.Server) {

	rd := redisCon.NewRedis(s.Cfg)

	inventoryRepo := repository.NewInventoryRepository(s.Db, rd)
	inventoryUsecase := useCase.NewInventoryUsecase(inventoryRepo)
	handler.NewInventoryHandler(inventoryUsecase)

	handler.NewInventoryGrpcHandler(inventoryUsecase)

	handler.NewInventoryQueueHandler(s.Cfg, inventoryUsecase)

}
