package itemShopController

import (
	"github.com/Montheankul-K/game-service/pkg/custom"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"github.com/Montheankul-K/game-service/pkg/itemShop/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type IItemShopControllerImpl struct {
	itemShopService itemShopService.IItemShopService
}

func NewItemShopController(itemShopService itemShopService.IItemShopService) IItemShopController {
	return &IItemShopControllerImpl{
		itemShopService: itemShopService,
	}
}

func (c *IItemShopControllerImpl) Listing(ctx echo.Context) error {
	itemFilter := new(itemShopModel.ItemFilter)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err := customEchoRequest.Bind(itemFilter); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	itemModelList, err := c.itemShopService.Listing(itemFilter)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, itemModelList)
}
