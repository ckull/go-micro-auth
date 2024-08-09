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

func (r *authRepository) init() *mongo.Collection {
	return r.db.Database("Auth").Collection("Users")
}

func (r *authRepository) findOneUserByEmail(email string) (*model.UserPassport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userPassport := new(model.UserPassport)

	authCollection := r.init()
	err := authCollection.FindOne(ctx, bson.M{"email": email}).Decode(&userPassport)
	if err != nil {
		return nil, err
	}

	return userPassport, err
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
