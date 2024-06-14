package itemManagingController

import (
	"fmt"
	"github.com/Montheankul-K/game-service/pkg/custom"
	itemManagingModel "github.com/Montheankul-K/game-service/pkg/itemManaging/model"
	itemManagingService "github.com/Montheankul-K/game-service/pkg/itemManaging/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type itemManagingController struct {
	itemManagingService itemManagingService.IItemManagingService
}

func NewItemManagingController(itemManagingService itemManagingService.IItemManagingService) IItemManagingController {
	return &itemManagingController{
		itemManagingService: itemManagingService,
	}
}

func (c *itemManagingController) Creating(ctx echo.Context) error {
	itemCreatingReq := new(itemManagingModel.ItemCreatingReq)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err := customEchoRequest.Bind(itemCreatingReq); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	item, err := c.itemManagingService.Creating(itemCreatingReq)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, item)
}

func (c *itemManagingController) Editing(ctx echo.Context) error {
	itemID, err := c.getItemID(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	itemEditingReq := new(itemManagingModel.ItemEditingReq)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err = customEchoRequest.Bind(itemEditingReq); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	item, err := c.itemManagingService.Editing(itemID, itemEditingReq)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, item)
}

func (c *itemManagingController) getItemID(ctx echo.Context) (uint64, error) {
	itemID := ctx.Param("itemID")
	itemIDUint64, err := strconv.ParseUint(itemID, 10, 64)
	if err != nil {
		return 0, err
	}
	return itemIDUint64, nil
}

func (c *itemManagingController) Archiving(ctx echo.Context) error {
	itemID, err := c.getItemID(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	if err = c.itemManagingService.Archiving(itemID); err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusNoContent, fmt.Sprintf("delete item id: %d successfully", itemID))
}
