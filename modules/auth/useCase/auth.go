package useCase

import (
	"errors"
	"fmt"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/modules/auth/repository"
	"go-auth/pkg/jwtAuth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthUsecase interface {
		RegisterByEmail(c echo.Context, cfg *config.Config, registerReq *model.RegisterReq) (*model.AccessToken, error)
		Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.AccessToken, error)
		Logout(c echo.Context, cfg *config.Config, logoutReq *model.LogoutReq) error
		ReloadToken(c echo.Context, cfg *config.Config, reloadReq *model.Token) (*string, error)
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

func (u *authUsecase) RegisterByEmail(c echo.Context, cfg *config.Config, registerReq *model.RegisterReq) (*model.AccessToken, error) {
	user, err := u.authRepository.FindOneUserByEmail(registerReq.Email)
	if err == nil && user != nil {
		return nil, model.ErrEmailAlreadyExists
	}

	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, model.ErrFailedToHashPassword
	}

	userPassport := &model.UserPassport{
		Email:         registerReq.Email,
		Password:      string(hashedPassword),
		OauthProvider: "email",
		Role:          "user",
	}

	if err := u.authRepository.AddUser(userPassport); err != nil {
		return nil, errors.New("failed to add user")
	}

	claims := &jwtAuth.Claims{
		UserId:   registerReq.Email,
		RoleCode: "user",
	}

	accessToken := u.authRepository.AccessToken(cfg, claims)

	refreshToken := u.authRepository.RefreshToken(cfg, claims)

	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(cfg.Jwt.RefreshTokenDuration) * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	c.SetCookie(refreshTokenCookie)

	return &model.AccessToken{
		AccessToken: accessToken,
	}, nil
}

func (u *authUsecase) Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.AccessToken, error) {
	user, err := u.authRepository.FindOneUserByEmail(loginReq.Email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, password is invalid")
	}

	claims := &jwtAuth.Claims{
		UserId:   user.Email,
		RoleCode: user.Role,
	}

	accessToken := u.authRepository.AccessToken(cfg, claims)

	refreshToken := u.authRepository.RefreshToken(cfg, claims)

	// Set refresh token as an HTTP-only cookie
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(cfg.Jwt.RefreshTokenDuration) * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	c.SetCookie(refreshTokenCookie)

	return &model.AccessToken{
		AccessToken: accessToken,
	}, nil

}

func (u *authUsecase) Logout(c echo.Context, cfg *config.Config, logoutReq *model.LogoutReq) error {
	// Parse the refresh token to extract the claims
	claims := &jwtAuth.AuthMapClaims{}
	token, err := jwt.ParseWithClaims(logoutReq.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.RefreshTokenSecret), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
	}

	// Get the expiration time from the token's claims
	expirationTime := claims.ExpiresAt.Time

	// Add the refresh token to the blacklist with the expiration time
	err = u.authRepository.AddBlacklistToken(logoutReq.RefreshToken, expirationTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to blacklist token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (u *authUsecase) ReloadToken(c echo.Context, cfg *config.Config, reloadReq *model.Token) (*string, error) {
	// Parse the access token
	accessClaims := &jwtAuth.AuthMapClaims{}
	accessToken, err := jwt.ParseWithClaims(reloadReq.AccessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.AccessTokenSecret), nil
	})

	// If the access token is valid, return it
	if err == nil && accessToken.Valid {
		return &reloadReq.AccessToken, nil // Access token is still valid
	}

	// If the access token is expired, check the refresh token
	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		refreshClaims := &jwtAuth.AuthMapClaims{}
		refreshToken, refreshTokenErr := jwt.ParseWithClaims(reloadReq.RefreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.RefreshTokenSecret), nil
		})

		if refreshTokenErr != nil || !refreshToken.Valid {
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired refresh token"})
		}

		if refreshTokenErr != nil && errors.Is(err, jwt.ErrTokenExpired) {
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "Expired refresh token"})
		}

		expirationTime := refreshClaims.ExpiresAt.Time

		if err := u.authRepository.AddBlacklistToken(reloadReq.RefreshToken, expirationTime); err != nil {
			return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "Add blacklist failed"})
		}

		// Generate new access token and refresh token
		newAccessToken := u.authRepository.AccessToken(cfg, refreshClaims.Claims)
		newRefreshToken := u.authRepository.RefreshToken(cfg, refreshClaims.Claims)

		// Set the new refresh token as an HTTP-only cookie
		refreshTokenCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    newRefreshToken,
			Expires:  time.Now().Add(time.Duration(cfg.Jwt.RefreshTokenDuration) * time.Hour),
			HttpOnly: true,
			Path:     "/",
		}

		c.SetCookie(refreshTokenCookie)

		return &newAccessToken, nil
	}

	// If the access token error is something else, return an error
	return nil, c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid access token"})
}
