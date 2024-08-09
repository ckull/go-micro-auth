package useCase

import (
	"errors"
	"fmt"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/modules/auth/repository"
	"go-auth/pkg/jwt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthUsecase interface {
	}

	authUsecase struct {
		authRepository repository.AuthRepository
	}
)

func NewAuthUsecase(authRepository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
	}
}

func (u *authUsecase) Register(c echo.Context, cfg *config.Config, registerReq *model.RegisterReq) (*model.AccessToken, error)  {
	userPassport, err := u.authRepository.findOneUserByEmail(registerReq.Email)
	if err == nil  && userPassport == nil {

	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	claims := &jwt.Claims{
		UserId:   registerReq.Email,
		RoleCode: "USER",
	}

	accessToken := jwt.GetTokens(c, cfg, u, claims)

	return accessToken, nil
}

func (u *authUsecase) Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.AccessToken, error) {
	userPassport, err := u.authRepository.findOneUserByEmail(loginReq.Email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userPassport.Password), []byte(loginReq.Password)); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, password is invalid")
	}

	claims := &jwt.Claims{
		UserId:   userPassport.Email,
		RoleCode: userPassport.RoleCode,
	}

	accessToken := u.authRepository.AccessToken(cfg, &claims)

	refreshToken := u.authRepository.RefreshToken(cfg, &claims)

	// Set refresh token as an HTTP-only cookie
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  cfg.Jwt.RefreshTokenDuration
		HttpOnly: true,
		Path:     "/",
	}

	c.SetCookie(refreshTokenCookie)

	return &AccessToken{
		access_token: accessToken
	}, nil

}
