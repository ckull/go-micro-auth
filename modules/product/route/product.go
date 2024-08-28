package route

import (
	"go-meechok/modules/product/handler"
	"go-meechok/modules/product/repository"
	"go-meechok/modules/product/service"
	"go-meechok/modules/product/useCase"
	"go-meechok/pkg/redisCon"
	"go-meechok/server/types"
)

func ProductRoute(s *types.Server) {
	rd := redisCon.NewRedis(s.Cfg)

	productRepo := repository.NewProductRepository(s.Db, rd)
	productUsecase := useCase.NewProductUsecase(productRepo)
	productService := service.NewProductQueueService(s.Cfg, productUsecase)
	handler.NewProductHandler(productUsecase)
	handler.NewProductGrpcHandler(productUsecase)

	handler.NewProductQueueHandler(s.Cfg, productUsecase, productService)
}
