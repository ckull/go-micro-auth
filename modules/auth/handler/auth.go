package handler

import (
	"errors"
	"fmt"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/modules/auth/useCase"
	"go-auth/pkg/cookieHelper"
	"go-auth/utils"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type (
	AuthHandler interface {
		RegisterByEmail(c echo.Context) error
		Login(c echo.Context) error
		Logout(c echo.Context) error
		RefreshToken(c echo.Context) error
	}

	authHandler struct {
		authUsecase useCase.AuthUsecase
		cfg         *config.Config
		validator   *validator.Validate
	}
)

func NewAuthHandler(authUsecase useCase.AuthUsecase, cfg *config.Config) AuthHandler {
	return &authHandler{
		authUsecase: authUsecase,
		cfg:         cfg,
		validator:   validator.New(),
	}
}

func (h *authHandler) RegisterByEmail(c echo.Context) error {
	var registerReq model.RegisterReq
	if err := c.Bind(&registerReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(registerReq); err != nil {
		validationErrors := utils.FormatValidationError(err)
		log.Printf("Error: Validate data failed: %s", err.Error())
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	fmt.Println("registerReq: ", registerReq)
	accessToken, err := h.authUsecase.RegisterByEmail(c, h.cfg, &registerReq)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEmailAlreadyExists):
			return c.JSON(http.StatusConflict, map[string]string{"error": "Email is already exist"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, accessToken)
}

func (h *authHandler) Login(c echo.Context) error {
	var loginReq model.LoginReq
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(loginReq); err != nil {
		validationErrors := utils.FormatValidationError(err)
		log.Printf("Error: Validate data failed: %s", err.Error())
		return c.JSON(http.StatusBadRequest, validationErrors)
	}

	accessToken, err := h.authUsecase.Login(c, h.cfg, &loginReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, accessToken)
}

func (h *authHandler) Logout(c echo.Context) error {
	var logoutReq *model.LogoutReq
	if err := c.Bind(logoutReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	err := h.authUsecase.Logout(c, h.cfg, logoutReq)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"error": "Successfully logout"})
}

func (h *authHandler) RefreshToken(c echo.Context) error {
	accessToken := c.Get("accessToken").(string)
	refreshToken := c.Get("refreshToken").(string)

	// Create a Token struct to pass to ReloadToken (assuming you have such a struct)
	reloadReq := &model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	newTokens, err := h.authUsecase.ReloadToken(c, h.cfg, reloadReq)

	cookie := cookieHelper.NewCookieHelper(c, h.cfg)

	cookie.SetRefreshToken(newTokens.RefreshToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, newTokens)

}
