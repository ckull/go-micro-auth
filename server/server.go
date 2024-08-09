package server

import (
	"context"
	"go-auth/config"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Server struct {
		App *echo.Echo
		Db  *mongo.Client
		Cfg *config.Config
	}
)

var (
	serverInstance Server
	once           sync.Once
)

func Start(ctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &Server{
		App: echo.New(),
		Db:  db,
		Cfg: cfg,
	}

	// CORS
	s.App.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
	}))

	s.App.Logger.Fatal(s.App.Start(":8080"))
}
