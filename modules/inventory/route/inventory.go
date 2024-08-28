package route

import (
	"go-meechok/modules/inventory/handler"
	"go-meechok/modules/inventory/repository"
	"go-meechok/modules/inventory/useCase"
	"go-meechok/pkg/redisCon"
	"go-meechok/server/types"
)

func AuthRoute(s *types.Server) {

	rd := redisCon.NewRedis(s.Cfg)

	inventoryRepo := repository.NewInventoryRepository(s.Db, rd)
	inventoryUsecase := useCase.NewInventoryUsecase(inventoryRepo)
	handler.NewInventoryHandler(inventoryUsecase)

	handler.NewInventoryGrpcHandler(inventoryUsecase)
	// oauthHandler := oauth.NewOAuthHandler(s.Cfg, authUsecase)

	// s.App.POST("/inventory/register/email", authHandler.RegisterByEmail)
	// s.App.POST("/inventory/login", authHandler.Login)
	// s.App.POST("/inventory/logout", authHandler.Logout)
	// s.App.POST("/inventory/refreshToken", authHandler.RefreshToken, middleware.JWTMiddleware())
	// s.App.GET("/inventory/users", authHandler.FindUserByUID, middleware.JWTMiddleware())
	// s.App.GET("/inventory/facebook/login", authHandler.FacebookLogin)
	// s.App.GET("/inventory/facebook/callback", authHandler.FacebookCallback)

	// s.App.GET("/auth/google/login", authHandler.FacebookLogin)
	// s.App.GET("/auth/google/callback", authHandler.FacebookCallback)

}
