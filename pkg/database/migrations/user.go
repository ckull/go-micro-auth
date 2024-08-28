package migrations

import (
	"context"
	"go-meechok/config"
	"go-meechok/pkg/database"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func userDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("User")
}

func UserMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// Users collection
	col := db.Collection("Users")

	// Create indexes for Users collection
	indexes, err := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"uid", 1}}, Options: options.Index().SetUnique(true)},
	})
	if err != nil {
		log.Fatalf("Error creating indexes for users collection: %v", err)
	}
	for _, index := range indexes {
		log.Printf("Created index: %s", index)
	}

	log.Println("User migrations completed successfully")
}
