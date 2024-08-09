package model

import "time"

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
		CredentialId string `json:"credential_id" form:"credential_id" validate:"required,max=64"`
	}

	UserPassport struct {
		Id       string `json:"_id"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	User struct {
		UID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		Email    string             `bson:"email" json:"email"`
		Password string             `bson:"password" json:"-"`
	}

	Token struct {
		UID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
		RefreshToken string             `bson:"refresh_token" json:"refresh_token"`
	}

	AccessToken struct {
		accessToken string `json:"access_token"`
	}
)
