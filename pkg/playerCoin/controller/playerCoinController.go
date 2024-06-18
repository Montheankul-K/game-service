package playerCoinController

import "github.com/labstack/echo/v4"

type IPlayerCoinController interface {
	CoinAdding(ctx echo.Context) error
	Showing(ctx echo.Context) error
}
