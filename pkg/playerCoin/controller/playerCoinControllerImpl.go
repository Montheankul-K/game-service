package playerCoinController

import (
	"github.com/Montheankul-K/game-service/pkg/custom"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	playerCoinService "github.com/Montheankul-K/game-service/pkg/playerCoin/service"
	"github.com/Montheankul-K/game-service/pkg/validation"
	"github.com/labstack/echo/v4"
	"net/http"
)

type playerCoinControllerImpl struct {
	playerCoinService playerCoinService.IPlayerCoinService
}

func NewPlayerCoinControllerImpl(playerCoinService playerCoinService.IPlayerCoinService) IPlayerCoinController {
	return &playerCoinControllerImpl{
		playerCoinService: playerCoinService,
	}
}

func (c *playerCoinControllerImpl) CoinAdding(ctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	coinAddingReq := new(playerCoinModel.CoinAddingReq)
	customEchoRequest := custom.NewCustomEchoRequest(ctx)
	if err = customEchoRequest.Bind(coinAddingReq); err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	coinAddingReq.PlayerID = playerID
	playerCoin, err := c.playerCoinService.CoinAdding(coinAddingReq)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusCreated, playerCoin)
}

func (c *playerCoinControllerImpl) Showing(ctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(ctx)
	if err != nil {
		return custom.Error(ctx, http.StatusBadRequest, err.Error())
	}

	playerCoinShowing := c.playerCoinService.Showing(playerID)
	return ctx.JSON(http.StatusOK, playerCoinShowing)
}
