package model

import (
	"time"
)

type (
	LoginReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
	}

	RefreshTokenReq struct {
		CredentialId string `json:"credential_id" form:"credential_id" validate:"required,max=64"`
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	RegisterReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		Role          string `json:"role"`
	}

	User struct {
		Email         string    `bson:"email" json:"email"`
		Password      string    `bson:"password" json:"password"`
		OauthProvider string    `bson:"oauth_provider" json:"oauth_provider"`
		Role          string    `bson:"role" json:"role"`
		CreatedAt     time.Time `bson:"created_at" json:"created_at"`
		UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
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
)
