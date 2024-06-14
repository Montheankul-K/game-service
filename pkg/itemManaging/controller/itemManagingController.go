package itemManagingController

import "github.com/labstack/echo/v4"

type IItemManagingController interface {
	Creating(ctx echo.Context) error
	Editing(ctx echo.Context) error
	Archiving(ctx echo.Context) error
}
