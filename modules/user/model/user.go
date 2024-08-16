package model

import (
	"time"
)

type (
	User struct {
		Uid       string    `bson:"uid" json:"uid"`
		Address   string    `bson:"address" json:"password"`
		FirstName string    `bson:"first_name" json:"first_name"`
		LastName  string    `bson:"last_name" json:"last_name"`
		Phone     string    `bson:"phone" json:"phone"`
		CreatedAt time.Time `bson:"created_at" json:"created_at"`
		UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
	}
)
