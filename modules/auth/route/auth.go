package route

import (
	"go-auth/middleware"
	"go-auth/modules/auth/handler"
	"go-auth/modules/auth/repository"
	"go-auth/modules/auth/useCase"
	"go-auth/pkg/oauth"
	"go-auth/server/types"
)

func AuthRoute(s *types.Server) {

	authRepo := repository.NewAuthRepository(s.Db)
	authUsecase := useCase.NewAuthUsecase(authRepo)
	authHandler := handler.NewAuthHandler(authUsecase, s.Cfg)

	oauthHandler := oauth.NewOAuthHandler(s, authUsecase)

	s.App.POST("auth/register/email", authHandler.RegisterByEmail)
	s.App.POST("auth/login", authHandler.Login)
	s.App.POST("auth/logout", authHandler.Logout)
	s.App.POST("auth/refreshToken", authHandler.RefreshToken, middleware.JWTMiddleware())

	s.App.GET("auth/facebook/login", oauthHandler.FacebookLogin)
	s.App.GET("auth/facebook/callback", oauthHandler.FacebookCallback)
}
