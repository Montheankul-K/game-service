package inventorycontroller

import "github.com/labstack/echo/v4"

type IInventoryController interface {
	Listing(ctx echo.Context) error
}
