package handler

import (
	"go-meechok/modules/user/useCase"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	uidParam := c.Param("uid")
	uid, err := primitive.ObjectIDFromHex(uidParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UID format"})
	}

	users, err := h.userUsecase.GetUserByUID(uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}
