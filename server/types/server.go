package types

import (
	"go-auth/config"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Server struct {
		App *echo.Echo
		Db  *mongo.Client
		Cfg *config.Config
	}
)
