package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Inventory struct {
		ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		ProductID primitive.ObjectID `bson:"product_id" json:"product_id"`
		Quantity  int                `bson:"quantity" json:"quantity"`
		Version   int64              `bson:"version" json:"version"`
		CreatedAt time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	}

	AddInventoryReq struct {
		ProductID primitive.ObjectID `json:"product_id`
		Quantity  int                `json:"product_id`
	}

	AddInventoryFailedEvent struct {
		ProductID string `json:"product_id`
	}
)
