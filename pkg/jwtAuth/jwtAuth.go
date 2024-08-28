package jwtAuth

import (
	"errors"
	"go-meechok/modules/auth/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	AuthFactory interface {
		SignToken() string
	}

	Claims struct {
		UserId string     `json:"user_id"`
		Role   model.Role `json:"role"`
	}

	AuthMapClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	authConcrete struct {
		Secret []byte
		Claims *AuthMapClaims `json:"claims"`
	}

	accessToken struct {
		*authConcrete
	}

	refreshToken struct {
		*authConcrete
	}

	apiKey struct {
		*authConcrete
	}
)

func (a *authConcrete) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims)
	ss, err := token.SignedString(a.Secret)
	if err != nil {

	}

	return ss
}

func now() time.Time {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {

	}
	return time.Now().In(loc)
}

// Note that: t is a second unit
func JwtTimeDurationMinute(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Minute * time.Duration(t)))
}

func JwtTimeDurationHour(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(now().Add(time.Hour * time.Duration(t)))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func NewAccessToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &accessToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    os.Getenv("JWT_ISSUER"),
					Subject:   "access-token",
					ExpiresAt: JwtTimeDurationMinute(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func NewRefreshToken(secret string, expiredAt int64, claims *Claims) AuthFactory {
	return &refreshToken{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: claims,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    os.Getenv("JWT_ISSUER"),
					Subject:   "refresh-token",
					ExpiresAt: JwtTimeDurationHour(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func NewApiKey(secret string, expiredAt int64) AuthFactory {
	return &apiKey{
		authConcrete: &authConcrete{
			Secret: []byte(secret),
			Claims: &AuthMapClaims{
				Claims: &Claims{},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    os.Getenv("JWT_ISSUER"),
					Subject:   "api-key",
					ExpiresAt: JwtTimeDurationHour(expiredAt),
					NotBefore: jwt.NewNumericDate(now()),
					IssuedAt:  jwt.NewNumericDate(now()),
				},
			},
		},
	}
}

func ParseToken(secret string, tokenString string) (*AuthMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthMapClaims{}, func(token *jwt.Token) (interface{}, error) {

		return []byte(secret), nil
	})

	if !token.Valid {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("error: token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("error: token is expired")
		} else {
			return nil, errors.New("error: token is invalid")
		}
	}

	if claims, ok := token.Claims.(*AuthMapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("error: claims type is invalid")
	}
}

// func GetTokens(c echo.Echo, cfg *config.Config, authUsecase *useCase.AuthUsecase, claims *Claims) string {
// 	accessToken := authUsecase.AuthRepository.AccessToken(cfg, &claims)

// 	refreshToken := authUsecase.AuthRepository.RefreshToken(cfg, &claims)

// 	refreshTokenCookie := &http.Cookie{
// 		Name:     "refresh_token",
// 		Value:    refreshToken,
// 		Expires:  JwtTimeDurationHour(cfg.Jwt.RefreshTokenDuration),
// 		HttpOnly: true,
// 		Path:     "/",
// 	}

// 	c.SetCookie(refreshTokenCookie)

// 	return accessToken

// }
