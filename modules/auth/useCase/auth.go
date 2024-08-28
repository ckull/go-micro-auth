package useCase

import (
	"errors"
	"fmt"
	"go-meechok/config"
	"go-meechok/modules/auth/model"
	"go-meechok/modules/auth/repository"
	"go-meechok/pkg/cookieHelper"
	"go-meechok/pkg/jwtAuth"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type (
	AuthUsecase interface {
		RegisterByEmail(c echo.Context, cfg *config.Config, registerReq *model.RegisterReq) (*model.AccessToken, error)
		Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.AccessToken, error)
		Logout(c echo.Context, cfg *config.Config, logoutReq *model.LogoutReq) error
		ReloadToken(c echo.Context, cfg *config.Config, reloadReq *model.Token) (*model.Token, error)
		FindOrRegisterFacebookUser(userInfo *model.FacebookUser) (*model.User, error)
		GenerateTokens(user *model.User, cfg *config.Config) *model.Token
		FindUserByUID(objectID primitive.ObjectID) (*model.User, error)
	}

	authUsecase struct {
		authRepository repository.AuthRepository
	}
)

func NewAuthUsecase(authRepository repository.AuthRepository) AuthUsecase {
	return &authUsecase{
		authRepository: authRepository,
	}
}

func (u *authUsecase) RegisterByEmail(c echo.Context, cfg *config.Config, registerReq *model.RegisterReq) (*model.AccessToken, error) {

	user, err := u.authRepository.FindOneUserByEmail(registerReq.Email)
	if err == nil && user != nil {
		return nil, model.ErrEmailAlreadyExists
	}

	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, model.ErrFailedToHashPassword
	}

	userPassport := &model.UserPassport{
		Email:         registerReq.Email,
		Password:      string(hashedPassword),
		OauthProvider: "email",
		Role:          "user",
	}

	newUser, err := u.authRepository.AddUser(userPassport)
	if err != nil {
		return nil, errors.New("failed to add user")
	}

	userId := newUser.ID.Hex()

	customerRole := &model.Role{
		Role:        "CUSTOMER",
		Permissions: []string{"VIEW_PRODUCT"},
	}

	claims := &jwtAuth.Claims{
		UserId: userId,
		Role:   *customerRole,
	}

	accessToken := u.authRepository.AccessToken(cfg, claims)

	refreshToken := u.authRepository.RefreshToken(cfg, claims)

	cookie := cookieHelper.NewCookieHelper(c, cfg)

	cookie.SetRefreshToken(refreshToken)

	return &model.AccessToken{
		AccessToken: accessToken,
	}, nil
}

func (u *authUsecase) FindUserByUID(objectID primitive.ObjectID) (*model.User, error) {
	user, error := u.authRepository.FindUserByUID(objectID)

	return user, error
}

func (u *authUsecase) FindOrRegisterFacebookUser(userInfo *model.FacebookUser) (*model.User, error) {
	facebookID := userInfo.OauthId
	user, err := u.authRepository.FindByProviderId(facebookID)

	if err != nil {
		return nil, err
	}

	if user == nil && err != nil {
		userPassport := &model.UserPassport{
			Email:         userInfo.Email,
			OauthProvider: "facebook",
			OauthId:       userInfo.OauthId,
		}

		if user, err := u.authRepository.AddUser(userPassport); err != nil {
			return user, err
		}
	}

	return user, err
}

func (u *authUsecase) Login(c echo.Context, cfg *config.Config, loginReq *model.LoginReq) (*model.AccessToken, error) {
	user, err := u.authRepository.FindOneUserByEmail(loginReq.Email)

	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("error, password is invalid")
	}

	userId := user.ID.Hex()

	claims := &jwtAuth.Claims{
		UserId: userId,
		Role:   user.Role,
	}

	accessToken := u.authRepository.AccessToken(cfg, claims)

	refreshToken := u.authRepository.RefreshToken(cfg, claims)

	cookie := cookieHelper.NewCookieHelper(c, cfg)

	cookie.SetRefreshToken(refreshToken)

	return &model.AccessToken{
		AccessToken: accessToken,
	}, nil

}

func (u *authUsecase) Logout(c echo.Context, cfg *config.Config, logoutReq *model.LogoutReq) error {
	// Parse the refresh token to extract the claims
	claims := &jwtAuth.AuthMapClaims{}
	token, err := jwt.ParseWithClaims(logoutReq.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.RefreshTokenSecret), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
	}

	if !token.Valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid or expired token"})
	}

	expirationTime := claims.ExpiresAt.Time

	err = u.authRepository.AddBlacklistToken(logoutReq.RefreshToken, expirationTime)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to blacklist token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (u *authUsecase) ReloadToken(c echo.Context, cfg *config.Config, reloadReq *model.Token) (*model.Token, error) {

	accessClaims := &jwtAuth.AuthMapClaims{}
	accessToken, err := jwt.ParseWithClaims(reloadReq.AccessToken, accessClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Jwt.AccessTokenSecret), nil
	})

	if err == nil && accessToken.Valid {
		return reloadReq, nil
	}

	if err != nil && errors.Is(err, jwt.ErrTokenExpired) {
		refreshClaims := &jwtAuth.AuthMapClaims{}
		refreshToken, refreshTokenErr := jwt.ParseWithClaims(reloadReq.RefreshToken, refreshClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.RefreshTokenSecret), nil
		})

		if refreshTokenErr != nil || !refreshToken.Valid {
			return nil, model.ErrInvalidRefreshToken
		}

		if refreshTokenErr != nil && errors.Is(err, jwt.ErrTokenExpired) {
			return nil, model.ErrExpiredRefreshToken
		}

		expirationTime := refreshClaims.ExpiresAt.Time

		if err := u.authRepository.AddBlacklistToken(reloadReq.RefreshToken, expirationTime); err != nil {
			return nil, model.ErrAddBlacklistTokenFailed
		}

		// Generate new access token and refresh token
		newAccessToken := u.authRepository.AccessToken(cfg, refreshClaims.Claims)
		newRefreshToken := u.authRepository.RefreshToken(cfg, refreshClaims.Claims)

		return &model.Token{
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		}, nil
	}

	// If the access token error is something else, return an error
	return nil, model.ErrInvalidAccessToken
}

func (u *authUsecase) GenerateTokens(user *model.User, cfg *config.Config) *model.Token {

	userId := user.ID.Hex()

	claims := &jwtAuth.Claims{
		UserId: userId,
		Role:   user.Role,
	}

	accessToken := u.authRepository.AccessToken(cfg, claims)

	refreshToken := u.authRepository.RefreshToken(cfg, claims)

	return &model.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
