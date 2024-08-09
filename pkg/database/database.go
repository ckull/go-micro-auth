package database

import (
	"context"
	"go-auth/config"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	once       sync.Once
	DbInstance *mongo.Client
)

func DbConn(pctx context.Context, cfg *config.Config) *mongo.Client {

	once.Do(func() {
		ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
		defer cancel()

		DbInstance, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.URI))

		if err != nil {
			log.Fatalf("Error: Conntect to database error: %s", err.Error())
		}

		if err := DbInstance.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalf("Error: Pinging to database error: %s", err.Error())
		}

	})

	return DbInstance

}
