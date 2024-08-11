package useCase

import (
	"go-auth/modules/user/model"
	"go-auth/modules/user/useCase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	UserHandler interface {
		GetUserByUID(c echo.Context) error
	}

	userHandler struct {
		userUsecase useCase.UserUsecase
	}
)

func NewUserHandler(userUsecase useCase.UserUsecase) UserHandler {
	return &userHandler{
		userUsecase: userUsecase,
	}
}

func (h *userHandler) GetUserByUID(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	users, err := h.userUsecase.GetUserByUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}
