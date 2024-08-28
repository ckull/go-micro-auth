package middleware

import (
	"go-meechok/pkg/jwtAuth"
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
	})
}
