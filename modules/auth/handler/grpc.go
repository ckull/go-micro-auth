package handler

import (
	"go-meechok/modules/auth/useCase"
)

type (
	authGrpcHandler struct {
		authUsecase useCase.AuthUsecase
	}
)

func NewAuthGrpcHandler(authUsecase useCase.AuthUsecase) *authGrpcHandler {
	return &authGrpcHandler{
		authUsecase: authUsecase,
	}
}
