package oauth2controller

import "github.com/labstack/echo/v4"

type IOAuth2Controller interface {
	PlayerLogin(ctx echo.Context) error
	AdminLogin(ctx echo.Context) error
	PlayerLoginCallback(ctx echo.Context) error
	AdminLoginCallback(ctx echo.Context) error
	Logout(ctx echo.Context) error
	PlayerAuthorizing(ctx echo.Context, next echo.HandlerFunc) error
	AdminAuthorizing(ctx echo.Context, next echo.HandlerFunc) error
}
