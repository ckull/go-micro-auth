package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	LoginReq struct {
		Email    string `json:"email" validate:"required,email,max=255"`
		Password string `json:"password" validate:"required,max=32"`
	}

	RefreshTokenReq struct {
		CredentialId string `json:"credential_id" form:"credential_id" validate:"required,max=64"`
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	RegisterReq struct {
		Email    string `json:"email" validate:"required,email,max=255"`
		Password string `json:"password" validate:"required,max=32"`
	}

	RegisterRes struct {
		AccessToken string `json:"access_token"`
	}

	LoginRes struct {
		Id           string    `json:"_id"`
		RoleCode     int       `json:"role_code"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	LogoutReq struct {
		RefreshToken string `json:"refresh_token"`
	}

	UserPassport struct {
		Email         string `json:"email"`
		Password      string `json:"password"`
		OauthProvider string `json:"oauth_provider"`
		OauthId       string `json:"oauth_id"`

		Role string `json:"role"`
	}

	User struct {
		ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Email         string             `bson:"email" json:"email"`
		Password      string             `bson:"password" json:"password,omitempty"`
		OauthProvider string             `bson:"oauth_provider" json:"oauth_provider"`
		OauthId       string             `bson:"oauth_id" json:"oauth_id,omitempty"`
		Role          string             `bson:"role" json:"role"`
		CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	}

	Blacklist struct {
		RefreshToken string    `bson:"refresh_token"`
		ExpiresAt    time.Time `bson:"expires_at"`
	}

	AccessToken struct {
		AccessToken string `json:"access_token"`
	}

	Token struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	FacebookUser struct {
		OauthId       string `json:"id"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		OauthProvider string `json:"oauth_provider"`
	}

	GoogleUser struct {
		OauthId       string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		OauthProvider string `json:"oauth_provider"`
	}
)
