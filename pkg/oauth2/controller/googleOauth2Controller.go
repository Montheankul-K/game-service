package oauth2controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Montheankul-K/game-service/config"
	adminModel "github.com/Montheankul-K/game-service/pkg/admin/model"
	"github.com/Montheankul-K/game-service/pkg/custom"
	oauth2Exception "github.com/Montheankul-K/game-service/pkg/oauth2/exception"
	oauth2Model "github.com/Montheankul-K/game-service/pkg/oauth2/model"
	oauth2Service "github.com/Montheankul-K/game-service/pkg/oauth2/service"
	playerModel "github.com/Montheankul-K/game-service/pkg/player/model"
	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type googleOAuth2Controller struct {
	oauth2Service oauth2Service.IOAuth2Service
	oauth2Cfg     *config.OAuth2
	logger        echo.Logger
}

var (
	playerGoogleOAuth2 *oauth2.Config
	adminGoogleOAuth2  *oauth2.Config
	once               sync.Once

	accessTokenCookieName  = "access"
	refreshTokenCookieName = "refresh"
	stateCookieName        = "state"

	letters = []byte("letter")
)

func NewGoogleOAuth2Controller(oauth2Service oauth2Service.IOAuth2Service, oauth2Cfg *config.OAuth2, logger echo.Logger) IOAuth2Controller {
	once.Do(func() {
		setGoogleOAuth2Config(oauth2Cfg)
	})
	return &googleOAuth2Controller{
		oauth2Service: oauth2Service,
		oauth2Cfg:     oauth2Cfg,
		logger:        logger,
	}
}

func setGoogleOAuth2Config(oauth2Cfg *config.OAuth2) {
	playerGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Cfg.ClientID,
		ClientSecret: oauth2Cfg.ClientSecret,
		RedirectURL:  oauth2Cfg.PlayerRedirectUrl,
		Scopes:       oauth2Cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Cfg.Endpoints.AuthUrl,
			TokenURL:      oauth2Cfg.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Cfg.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}

	adminGoogleOAuth2 = &oauth2.Config{
		ClientID:     oauth2Cfg.ClientID,
		ClientSecret: oauth2Cfg.ClientSecret,
		RedirectURL:  oauth2Cfg.AdminRedirectUrl,
		Scopes:       oauth2Cfg.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:       oauth2Cfg.Endpoints.AuthUrl,
			TokenURL:      oauth2Cfg.Endpoints.TokenUrl,
			DeviceAuthURL: oauth2Cfg.Endpoints.DeviceAuthUrl,
			AuthStyle:     oauth2.AuthStyleInParams,
		},
	}
}

func (c *googleOAuth2Controller) PlayerLogin(ctx echo.Context) error {
	state := c.randomState()
	c.setCookie(ctx, stateCookieName, state)
	return ctx.Redirect(http.StatusFound, playerGoogleOAuth2.AuthCodeURL(state))
}

func (c *googleOAuth2Controller) AdminLogin(ctx echo.Context) error {
	state := c.randomState()
	c.setCookie(ctx, stateCookieName, state)
	return ctx.Redirect(http.StatusFound, adminGoogleOAuth2.AuthCodeURL(state))
}

func (c *googleOAuth2Controller) setCookie(ctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}
	ctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) removeCookie(ctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	ctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) setSameSiteCookie(ctx echo.Context, name, value string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	ctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) removeSameSiteCookie(ctx echo.Context, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
	ctx.SetCookie(cookie)
}

func (c *googleOAuth2Controller) PlayerLoginCallback(ctx echo.Context) error {
	contxt := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(ctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("failed to validate callback: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	unauthorizedError := func() error {
		return &oauth2Exception.Unauthorized{}
	}
	token, err := playerGoogleOAuth2.Exchange(contxt, ctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("failed to exchange token: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	client := playerGoogleOAuth2.Client(contxt, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("failed to get user info: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	playerCreatingReq := &playerModel.PlayerCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	oauthProcessingError := func() error {
		return &oauth2Exception.OAuth2Processing{}
	}
	if err = c.oauth2Service.PlayerAccountCreating(playerCreatingReq); err != nil {
		c.logger.Errorf("failed to create account: %s", err.Error())
		return custom.Error(ctx, http.StatusInternalServerError, oauthProcessingError().Error())
	}

	c.setSameSiteCookie(ctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(ctx, refreshTokenCookieName, token.RefreshToken)
	return ctx.JSON(http.StatusOK, &oauth2Model.LoginResponse{
		Message: "login success",
	})
}

func (c *googleOAuth2Controller) AdminLoginCallback(ctx echo.Context) error {
	contxt := context.Background()
	if err := retry.Do(func() error {
		return c.callbackValidating(ctx)
	}, retry.Attempts(3), retry.Delay(3*time.Second)); err != nil {
		c.logger.Errorf("failed to validate callback: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	unauthorizedError := func() error {
		return &oauth2Exception.Unauthorized{}
	}
	token, err := adminGoogleOAuth2.Exchange(contxt, ctx.QueryParam("code"))
	if err != nil {
		c.logger.Errorf("failed to exchange token: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	client := adminGoogleOAuth2.Client(contxt, token)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		c.logger.Errorf("failed to get user info: %s", err.Error())
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	adminCreatingReq := &adminModel.AdminCreatingReq{
		ID:     userInfo.ID,
		Email:  userInfo.Email,
		Name:   userInfo.Name,
		Avatar: userInfo.Picture,
	}

	oauthProcessingError := func() error {
		return &oauth2Exception.OAuth2Processing{}
	}
	if err = c.oauth2Service.AdminAccountCreating(adminCreatingReq); err != nil {
		c.logger.Errorf("failed to create account: %s", err.Error())
		return custom.Error(ctx, http.StatusInternalServerError, oauthProcessingError().Error())
	}

	c.setSameSiteCookie(ctx, accessTokenCookieName, token.AccessToken)
	c.setSameSiteCookie(ctx, refreshTokenCookieName, token.RefreshToken)
	return ctx.JSON(http.StatusOK, &oauth2Model.LoginResponse{
		Message: "login success",
	})
}

func (c *googleOAuth2Controller) Logout(ctx echo.Context) error {
	logoutError := func() error {
		return &oauth2Exception.Logout{}
	}
	accessToken, err := ctx.Cookie(accessTokenCookieName)
	if err != nil {
		c.logger.Errorf("failed to get access token: %s", err.Error())
		return custom.Error(ctx, http.StatusBadRequest, logoutError().Error())
	}

	if err = c.revokeToken(accessToken.Value); err != nil {
		c.logger.Errorf("failed to revoke token: %s", err.Error())
		return custom.Error(ctx, http.StatusInternalServerError, logoutError().Error())
	}

	c.removeSameSiteCookie(ctx, accessTokenCookieName)
	c.removeSameSiteCookie(ctx, refreshTokenCookieName)
	return ctx.JSON(http.StatusOK, &oauth2Model.LogoutResponse{
		Message: "logout success",
	})
}

func (c *googleOAuth2Controller) revokeToken(accessToken string) error {
	revokeUrl := fmt.Sprintf("%s?token=%s", c.oauth2Cfg.RevokeUrl, accessToken)
	res, err := http.Post(revokeUrl, "application/x-www-form-urlencoded", nil)
	if err != nil {
		c.logger.Errorf("failed to revoke token: %s", err.Error())
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *googleOAuth2Controller) getUserInfo(client *http.Client) (*oauth2Model.UserInfo, error) {
	res, err := client.Get(c.oauth2Cfg.UserInfoUrl)
	if err != nil {
		c.logger.Errorf("failed to get user info: %s", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	userInfoBytes, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Errorf("failed to read user info: %s", err.Error())
		return nil, err
	}

	userInfo := new(oauth2Model.UserInfo)
	if err = json.Unmarshal(userInfoBytes, &userInfo); err != nil {
		c.logger.Errorf("failed to unmarshal user info: %s", err.Error())
		return nil, err
	}
	return userInfo, nil
}

func (c *googleOAuth2Controller) callbackValidating(ctx echo.Context) error {
	state := ctx.QueryParam("state")
	stateFormCookie, err := ctx.Cookie(stateCookieName)
	if err != nil {
		c.logger.Errorf("failed to get state from cookie: %s", err.Error())
		return &oauth2Exception.Unauthorized{}
	}

	if state != stateFormCookie.Value {
		c.logger.Errorf("invalid state from cookie: %s", state)
		return &oauth2Exception.Unauthorized{}
	}

	c.removeCookie(ctx, stateCookieName)
	return nil
}

func (c *googleOAuth2Controller) randomState() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
