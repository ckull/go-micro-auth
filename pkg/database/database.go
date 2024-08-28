package database

import (
	"context"
	"go-meechok/config"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

		var err error
		DbInstance, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.Db.URI))

		if err != nil {
			log.Fatalf("Error: Connect to database error: %s", err.Error())
		}

		// Ping to ensure that the connection is established
		if err := DbInstance.Ping(ctx, readpref.Primary()); err != nil {
			log.Fatalf("Error: Pinging to database error: %s", err.Error())
		}

		log.Println("Successfully connected to the database.")

		// Create indexes
		// if err := createIndexes(ctx, DbInstance); err != nil {
		// 	log.Fatalf("Error: Create indexes: %s", err.Error())
		// }

		log.Println("Indexes created successfully.")
	})

	return DbInstance
}

func createIndexes(pctx context.Context, db *mongo.Client) error {
	_, err := db.Database("Auth").Collection("Users").Indexes().CreateOne(
		pctx, mongo.IndexModel{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.M{"oauth_provider": "email"}),
		},
	)

	if err != nil {
		return err
	}

	log.Println("Unique index on 'email' for 'Users' collection created successfully.")
	return nil
}
