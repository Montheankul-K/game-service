package oauth2controller

import (
	"context"
	"github.com/Montheankul-K/game-service/pkg/custom"
	oauth2Exception "github.com/Montheankul-K/game-service/pkg/oauth2/exception"
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"net/http"
)

func (c *googleOAuth2Controller) PlayerAuthorizing(ctx echo.Context, next echo.HandlerFunc) error {
	contxt := context.Background()
	tokenSource, err := c.getTokenSource(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	if !tokenSource.Valid() {
		tokenSource, err = c.playerTokenRefreshing(ctx, tokenSource)
		if err != nil {
			return custom.Error(ctx, http.StatusUnauthorized, err.Error())
		}
	}

	client := playerGoogleOAuth2.Client(contxt, tokenSource)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	unauthorizedError := func() error {
		return &oauth2Exception.Unauthorized{}
	}
	if !c.oauth2Service.IsPlayer(userInfo.ID) {
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	ctx.Set("playerID", userInfo.ID)
	return next(ctx)
}

func (c *googleOAuth2Controller) AdminAuthorizing(ctx echo.Context, next echo.HandlerFunc) error {
	contxt := context.Background()
	tokenSource, err := c.getTokenSource(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	if !tokenSource.Valid() {
		tokenSource, err = c.adminTokenRefreshing(ctx, tokenSource)
		if err != nil {
			return custom.Error(ctx, http.StatusUnauthorized, err.Error())
		}
	}

	client := adminGoogleOAuth2.Client(contxt, tokenSource)
	userInfo, err := c.getUserInfo(client)
	if err != nil {
		return custom.Error(ctx, http.StatusUnauthorized, err.Error())
	}

	unauthorizedError := func() error {
		return &oauth2Exception.Unauthorized{}
	}
	if !c.oauth2Service.IsAdmin(userInfo.ID) {
		return custom.Error(ctx, http.StatusUnauthorized, unauthorizedError().Error())
	}

	ctx.Set("adminID", userInfo.ID)
	return next(ctx)
}

func (c *googleOAuth2Controller) playerTokenRefreshing(ctx echo.Context, token *oauth2.Token) (*oauth2.Token, error) {
	contxt := context.Background()
	updatedToken, err := playerGoogleOAuth2.TokenSource(contxt, token).Token()
	if err != nil {
		return nil, &oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(ctx, accessTokenCookieName, updatedToken.AccessToken)
	c.setSameSiteCookie(ctx, refreshTokenCookieName, updatedToken.RefreshToken)
	return updatedToken, nil
}

func (c *googleOAuth2Controller) adminTokenRefreshing(ctx echo.Context, token *oauth2.Token) (*oauth2.Token, error) {
	contxt := context.Background()
	updatedToken, err := adminGoogleOAuth2.TokenSource(contxt, token).Token()
	if err != nil {
		return nil, &oauth2Exception.Unauthorized{}
	}

	c.setSameSiteCookie(ctx, accessTokenCookieName, updatedToken.AccessToken)
	c.setSameSiteCookie(ctx, refreshTokenCookieName, updatedToken.RefreshToken)
	return updatedToken, nil
}

func (c *googleOAuth2Controller) getTokenSource(ctx echo.Context) (*oauth2.Token, error) {
	accessToken, err := ctx.Cookie(accessTokenCookieName)
	if err != nil {
		return nil, &oauth2Exception.Unauthorized{}
	}

	refreshToken, err := ctx.Cookie(refreshTokenCookieName)
	if err != nil {
		return nil, &oauth2Exception.Unauthorized{}
	}

	return &oauth2.Token{
		AccessToken:  accessToken.Value,
		RefreshToken: refreshToken.Value,
	}, nil
}
