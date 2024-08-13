package migrations

import (
	"context"
	"go-auth/config"
	"go-auth/pkg/database"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("Auth")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	// Users collection
	col := db.Collection("Users")

	// Create indexes for Users collection
	indexes, err := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{{"email", 1}}, Options: options.Index().SetUnique(true)},
		// {Keys: bson.D{{"oauth_provider", 1}}},
		// {Keys: bson.D{{"role", 1}}},
	})
	if err != nil {
		log.Fatalf("Error creating indexes for users collection: %v", err)
	}
	for _, index := range indexes {
		log.Printf("Created index: %s", index)
	}

	// Blacklist collection
	col = db.Collection("Blacklists")

	// Create indexes for Blacklist collection
	indexes, err = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		// {Keys: bson.D{{"refresh_token", 1}}},
		// {Keys: bson.D{{"expires_at", 1}}},
	})
	if err != nil {
		log.Fatalf("Error creating indexes for blacklist collection: %v", err)
	}
	for _, index := range indexes {
		log.Printf("Created index: %s", index)
	}

	log.Println("User and Blacklist migrations completed successfully")
}
