package server

import (
	"context"
	"go-auth/config"
	auth "go-auth/modules/auth/route"
	user "go-auth/modules/user/route"
	"go-auth/server/types"

	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	serverInstance types.Server
	once           sync.Once
)

func Start(ctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &types.Server{
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

	auth.AuthRoute(s)
	user.UserRoute(s)
	s.App.Logger.Fatal(s.App.Start(":8080"))
}
