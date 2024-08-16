package middleware

import (
	"go-auth/pkg/jwtAuth"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echoJwt.WithConfig(echoJwt.Config{
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")), // Replace with your actual secret key
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtAuth.AuthMapClaims)
		},
		SuccessHandler: func(c echo.Context) {
			// Access token is valid, extract it from the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader != "" {
				c.Set("accessToken", authHeader)
			}
			// Extract the refresh token from the cookie
			refreshToken, err := c.Cookie("refresh_token")
			if err == nil {
				c.Set("refreshToken", refreshToken.Value)
			}
		},
	})
}
