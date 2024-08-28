package route

import (
	"go-meechok/middleware"
	"go-meechok/modules/auth/handler"
	"go-meechok/modules/auth/repository"
	"go-meechok/modules/auth/useCase"
	"go-meechok/pkg/redisCon"
	"go-meechok/server/types"
)

func AuthRoute(s *types.Server) {

	rd := redisCon.NewRedis(s.Cfg)

	authRepo := repository.NewAuthRepository(s.Db, rd)
	authUsecase := useCase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase, s.Cfg)

	handler.NewAuthGrpcHandler(authUsecase)
	// oauthHandler := oauth.NewOAuthHandler(s.Cfg, authUsecase)

	s.App.POST("/auth/register/email", authHandler.RegisterByEmail)
	s.App.POST("/auth/login", authHandler.Login)
	s.App.POST("/auth/logout", authHandler.Logout)
	s.App.POST("/auth/refreshToken", authHandler.RefreshToken, middleware.JWTMiddleware())
	s.App.GET("/auth/users", authHandler.FindUserByUID, middleware.JWTMiddleware())
	s.App.GET("/auth/facebook/login", authHandler.FacebookLogin)
	s.App.GET("/auth/facebook/callback", authHandler.FacebookCallback)

	s.App.GET("/auth/google/login", authHandler.FacebookLogin)
	s.App.GET("/auth/google/callback", authHandler.FacebookCallback)

}
