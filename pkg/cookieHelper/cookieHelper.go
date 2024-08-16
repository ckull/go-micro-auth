package cookieHelper

import (
	"go-auth/config"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Cookie interface {
		SetRefreshToken(refreshToken string)
	}

	cookie struct {
		Cfg     *config.Config
		Context echo.Context
	}
)

func NewCookieHelper(c echo.Context, cfg *config.Config) Cookie {
	return &cookie{
		Context: c,
		Cfg:     cfg,
	}
}

func (c *cookie) SetRefreshToken(refreshToken string) {
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Duration(c.Cfg.RefreshTokenDuration) * time.Hour),
		HttpOnly: true,
		Path:     "/",
	}

	c.Context.SetCookie(refreshTokenCookie)
}
