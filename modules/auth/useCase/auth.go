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

func (u *authUsecase) Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.LoginRes, error) {
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

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Minute * 15), // Set the expiration time
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(accessTokenCookie)

	// Set refresh token as an HTTP-only cookie
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7), // Set the expiration time for the refresh token
		HttpOnly: true,
		Path:     "/",
	}

	c.SetCookie(refreshTokenCookie)

}
