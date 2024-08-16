package repository

import (
	"context"
	"go-auth/config"
	"go-auth/modules/auth/model"
	"go-auth/pkg/jwtAuth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	AuthRepository interface {
		userCollection() *mongo.Collection
		blacklistCollection() *mongo.Collection
		FindOneUserByEmail(email string) (*model.User, error)
		AccessToken(cfg *config.Config, claims *jwtAuth.Claims) string
		RefreshToken(cfg *config.Config, claims *jwtAuth.Claims) string
		AddUser(userPassport *model.UserPassport) (*model.User, error)
		AddBlacklistToken(refreshToken string, expiration time.Time) error
		IsBlacklistExist(refreshToken string) (bool, error)
		FindByProviderId(id string) (*model.User, error)
		FindUserByUID(objectID primitive.ObjectID) (*model.User, error)
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

func (r *authRepository) userCollection() *mongo.Collection {
	return r.db.Database("Auth").Collection("Users")
}

func (r *authRepository) blacklistCollection() *mongo.Collection {
	return r.db.Database("Auth").Collection("Blacklists")
}

func (r *authRepository) FindByProviderId(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.userCollection()

	var user model.User
	err := collection.FindOne(ctx, bson.M{"oauth_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err

}

func (r *authRepository) AddUser(userPassport *model.UserPassport) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := &model.User{
		Email:         userPassport.Email,
		Password:      userPassport.Password,
		OauthProvider: userPassport.OauthProvider,
		OauthId:       userPassport.OauthId,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if userPassport.OauthId != "" {
		newUser.OauthId = userPassport.OauthId
	}

	collection := r.userCollection()

	_, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return newUser, err
}

func (r *authRepository) FindOneUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user := new(model.User)

	collection := r.userCollection()
	filter := bson.M{"email": email}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *authRepository) FindUserByUID(objectID primitive.ObjectID) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.userCollection()
	var user *model.User
	filter := bson.M{"_id": objectID}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *authRepository) AddBlacklistToken(refreshToken string, expiration time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.blacklistCollection()

	blacklist := bson.M{
		"refresh_token": refreshToken,
		"expires_at":    expiration,
	}

	_, err := collection.InsertOne(ctx, blacklist)
	return err
}

func (r *authRepository) IsBlacklistExist(refreshToken string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.blacklistCollection()
	filter := bson.M{
		"refresh_token": refreshToken,
		"expires_at":    bson.M{"$gt": time.Now()},
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *authRepository) AccessToken(cfg *config.Config, claims *jwtAuth.Claims) string {
	return jwtAuth.NewAccessToken(cfg.Jwt.AccessTokenSecret, cfg.Jwt.AccessTokenDuration, &jwtAuth.Claims{
		UserId:   claims.UserId,
		RoleCode: claims.RoleCode,
	}).SignToken()
}

func (r *authRepository) RefreshToken(cfg *config.Config, claims *jwtAuth.Claims) string {
	return jwtAuth.NewRefreshToken(cfg.Jwt.RefreshTokenSecret, cfg.Jwt.RefreshTokenDuration, &jwtAuth.Claims{
		UserId:   claims.UserId,
		RoleCode: claims.RoleCode,
	}).SignToken()
}
