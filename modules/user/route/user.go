package route

import (
	"go-auth/modules/user/handler"
	"go-auth/modules/user/middleware"
	"go-auth/modules/user/repository"
	"go-auth/modules/user/useCase"
	"go-auth/server/types"
)

func UserRoute(s *types.Server) {

	userRepo := repository.NewUserRepository(s.Db)
	userUsecase := useCase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	s.App.GET("/users", userHandler.GetUserByUID, middleware.JWTMiddleware())
}
