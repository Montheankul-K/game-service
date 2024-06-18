package inventorycontroller

import (
	"github.com/Montheankul-K/game-service/pkg/custom"
	inventoryService "github.com/Montheankul-K/game-service/pkg/inventory/service"
	"github.com/Montheankul-K/game-service/pkg/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type inventoryControllerImpl struct {
	inventoryService inventoryService.IInventoryService
	logger           echo.Logger
}

func NewInventoryControllerImpl(inventoryService inventoryService.IInventoryService, logger echo.Logger) IInventoryController {
	return &inventoryControllerImpl{
		inventoryService: inventoryService,
		logger:           logger,
	}
}

func (c *inventoryControllerImpl) Listing(ctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(ctx)
	if err != nil {
		c.logger.Errorf("failed to getting player id: %s", err.Error())
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	inventoryListing, err := c.inventoryService.Listing(playerID)
	if err != nil {
		return custom.Error(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, inventoryListing)
}
