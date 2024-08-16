package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/modules/auth/useCase"
	"go-auth/pkg/cookieHelper"
	"go-auth/pkg/jwtAuth"
	"go-auth/utils"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	AuthHandler interface {
		RegisterByEmail(c echo.Context) error
		Login(c echo.Context) error
		Logout(c echo.Context) error
		RefreshToken(c echo.Context) error
		FacebookLogin(c echo.Context) error
		FacebookCallback(c echo.Context) error
		FindUserByUID(c echo.Context) error
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

func (h *authHandler) FacebookLogin(c echo.Context) error {
	fmt.Println("facebook: ", h.cfg.Facebook)

	redirectUrl := h.cfg.Facebook.RedirectURL
	fmt.Println("redirect url: ", redirectUrl)

	fmt.Println("echo context: ", c)

	return c.Redirect(http.StatusTemporaryRedirect, h.cfg.Facebook.RedirectURL)
}

func (h *authHandler) FindUserByUID(c echo.Context) error {
	userJwt, ok := c.Get("user").(*jwt.Token)

	if !ok {
		return errors.New("JWT token missing or invalid")
	}

	claims, ok := userJwt.Claims.(*jwtAuth.AuthMapClaims)
	if !ok {
		return errors.New("JWT token missing or invalid")
	}

	fmt.Println(claims.UserId)

	uid, err := primitive.ObjectIDFromHex(claims.UserId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UID format"})
	}

	user, err := h.authUsecase.FindUserByUID(uid)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid UID format"})
	}

	return c.JSON(http.StatusOK, user)
}

func (h *authHandler) GoogleLogin(c echo.Context) error {
	redirectUrl := h.cfg.Google.AuthCodeURL("state")

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (a *authHandler) getFacebookUserInfo(accessToken string) (*model.FacebookUser, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo model.FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (h *authHandler) FacebookCallback(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	code := c.QueryParam("code")

	token, err := h.cfg.Facebook.Exchange(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
	}

	userInfo, err := h.getFacebookUserInfo(token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve user info"})
	}

	user, err := h.authUsecase.FindOrRegisterFacebookUser(userInfo)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get user info: "+err.Error())
	}

	tokens := h.authUsecase.GenerateTokens(user, h.cfg)

	cookie := cookieHelper.NewCookieHelper(c, h.cfg)

	cookie.SetRefreshToken(tokens.RefreshToken)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}
