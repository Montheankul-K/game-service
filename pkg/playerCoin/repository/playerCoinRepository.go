package playerCoinRepository

import (
	"github.com/Montheankul-K/game-service/entities"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	"gorm.io/gorm"
)

type IPlayerCoinRepository interface {
	CoinAdding(tx *gorm.DB, playerCoinEntity *entities.PlayerCoin) (*entities.PlayerCoin, error)
	Showing(playerID string) (*playerCoinModel.PlayerCoinShowing, error)
}
