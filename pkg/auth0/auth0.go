package auth0

import (
	"context"
	"encoding/json"
	"go-auth/modules/auth/model"
	"go-auth/modules/auth/useCase"
	"go-auth/pkg/cookieHelper"
	"go-auth/server/types"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type (
	Auth0Provider interface {
	}

	auth0Provider struct {
		Server      *types.Server
		AuthUsecase useCase.AuthUsecase
	}
)

func NewAuth0Provider(server *types.Server, authUsecase useCase.AuthUsecase) Auth0Provider {
	return &auth0Provider{
		Server:      server,
		AuthUsecase: authUsecase,
	}
}

func (a *auth0Provider) facebookLogin(c echo.Context) error {
	redirectUrl := a.Server.Cfg.Facebook.AuthCodeURL("state")

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (a *auth0Provider) googleLogin(c echo.Context) error {
	redirectUrl := a.Server.Cfg.Google.AuthCodeURL("state")

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (a *auth0Provider) getFacebookUserInfo(accessToken string) (*model.FacebookUser, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://graph.facebook.com/me?fields=id,name,email", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo model.FacebookUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func (a *auth0Provider) facebookCallback(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	code := c.QueryParam("code")

	token, err := a.Server.Cfg.Facebook.Exchange(ctx, code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
	}

	userInfo, err := a.getFacebookUserInfo(token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve user info"})
	}

	user, err := a.AuthUsecase.FindOrRegisterFacebookUser(userInfo)

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get user info: "+err.Error())
	}

	tokens := a.AuthUsecase.GenerateTokens(user, a.Server.Cfg)

	cookie := cookieHelper.NewCookieHelper(c, a.Server.Cfg)

	cookie.SetRefreshToken(tokens.RefreshToken)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}
