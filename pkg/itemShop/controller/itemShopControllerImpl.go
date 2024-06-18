package itemShopController

import (
	"github.com/Montheankul-K/game-service/pkg/custom"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"github.com/Montheankul-K/game-service/pkg/itemShop/service"
	"github.com/Montheankul-K/game-service/pkg/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type itemShopControllerImpl struct {
	itemShopService itemShopService.IItemShopService
}

func NewItemShopController(itemShopService itemShopService.IItemShopService) IItemShopController {
	return &itemShopControllerImpl{
		itemShopService: itemShopService,
	}
}

func (c *itemShopControllerImpl) Listing(ctx echo.Context) error {
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

func (c *itemShopControllerImpl) Buying(ctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	buyingReq := new(itemShopModel.BuyingReq)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err = customEchoRequest.Bind(buyingReq); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}
	buyingReq.PlayerID = playerID

	playerCoin, err := c.itemShopService.Buying(buyingReq)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, playerCoin)
}

func (c *itemShopControllerImpl) Selling(ctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	sellingReq := new(itemShopModel.SellingReq)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err = customEchoRequest.Bind(sellingReq); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}
	sellingReq.PlayerID = playerID

	playerCoin, err := c.itemShopService.Selling(sellingReq)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, playerCoin)
}
