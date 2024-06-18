package itemShopController

import "github.com/labstack/echo/v4"

type IItemShopController interface {
	Listing(ctx echo.Context) error
	Buying(ctx echo.Context) error
	Selling(ctx echo.Context) error
}
