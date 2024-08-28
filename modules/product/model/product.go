package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Product struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty`
		Name        string             `bson:"name"	json:"name"`
		Description string             `bson:"description" json:"description"`
		Price       float64            `bson:"price" json:"price"`
		Categories  []string           `bson:"categories" json:"categories"`
		IsActive    bool               `bson:"is_active" json:"is_active"`
		SellerID    primitive.ObjectID `bson:"seller_id" json:"seller_id"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}

	CreateProduct struct {
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Price       float64            `json:"price"`
		Categories  []string           `json:"categories"`
		IsActive    bool               `json:"is_active"`
		SellerID    primitive.ObjectID `json:"seller_id"`
		Quantity    int                `json:"quantity"`
	}

	ProductCreatedPayload struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
)
