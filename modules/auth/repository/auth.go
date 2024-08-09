package repository

import (
	"context"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/pkg/jwt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepository interface {
	}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepository {
	return &authRepository{
		db,
	}
}

func (r *authRepository) usersCollection() *mongo.Collection {
	return r.db.Database("Auth").Collection("Users")
}

func (r *authRepository) tokensCollection() *mongo.Collection {
	return r.db.Database("Auth").Collection("Tokens")
}

func (r *authRepository) findOneUserByEmail(email string) (*model.UserPassport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userPassport := new(model.UserPassport)

	collection := r.usersCollection()
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&userPassport)
	if err != nil {
		return nil, err
	}

	return userPassport, err
}

func (r *authRepository) updateToken(uid string, refreshToken string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.tokensCollection()

	filter := bson.M{"_id": uid}
	update := bson.M{"$set": bson.M{"refresh_token": refreshToken}}
	result, err := refreshTokenCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	token := new(model.Token)
	updated := result.Decode(&token)

	return updated, nil
}

func (r *authRepository) AccessToken(cfg *config.Config, claims *jwt.Claims) string {
	return jwt.NewAccessToken(cfg.Jwt.AccessTokenSecret, cfg.Jwt.AccessTokenDuration, &jwt.Claims{
		UserId:   claims.UserId,
		RoleCode: claims.RoleCode,
	}).SignToken()
}

func (r *authRepository) RefreshToken(cfg *config.Config, claims *jwt.Claims) string {
	return jwt.NewRefreshToken(cfg.Jwt.RefreshTokenSecret, cfg.Jwt.RefreshTokenDuration, &jwt.Claims{
		UserId:   claims.UserId,
		RoleCode: claims.RoleCode,
	}).SignToken()
}
