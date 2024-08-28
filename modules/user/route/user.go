package route

import (
	"go-meechok/modules/user/handler"
	"go-meechok/modules/user/middleware"
	"go-meechok/modules/user/repository"
	"go-meechok/modules/user/useCase"
	"go-meechok/server/types"
)

func UserRoute(s *types.Server) {

	userRepo := repository.NewUserRepository(s.Db)
	userUsecase := useCase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	s.App.GET("/users", userHandler.GetUserByUID, middleware.JWTMiddleware())
}
