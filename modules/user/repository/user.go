package repository

import (
	"context"
	"go-auth/modules/user/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	UserRepository interface {
		userCollection() *mongo.Collection
		GetUserByUID(uid primitive.ObjectID) (model.User, error)
	}

	userRepository struct {
		db *mongo.Client
	}
)

func NewUserRepository(db *mongo.Client) UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) userCollection() *mongo.Collection {
	return r.db.Database("User").Collection("Users")
}

func (r *userRepository) GetUserByUID(uid primitive.ObjectID) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	userCollection := r.userCollection()
	err := userCollection.FindOne(ctx, bson.M{"uid": uid}).Decode(&user)

	return user, err
}
