package validation

import (
	adminException "github.com/Montheankul-K/game-service/pkg/admin/exception"
	playerException "github.com/Montheankul-K/game-service/pkg/player/exception"
	"github.com/labstack/echo/v4"
)

func AdminIDGetting(ctx echo.Context) (string, error) {
	if adminID, ok := ctx.Get("adminID").(string); !ok || adminID == "" {
		return "", &adminException.AdminNotFound{
			AdminID: "Unknown",
		}
	} else {
		return adminID, nil
	}
}

func PlayerIDGetting(ctx echo.Context) (string, error) {
	if playerID, ok := ctx.Get("playerID").(string); !ok || playerID == "" {
		return "", &playerException.PlayerNotFound{
			PlayerID: "Unknown",
		}
	} else {
		return playerID, nil
	}
}
