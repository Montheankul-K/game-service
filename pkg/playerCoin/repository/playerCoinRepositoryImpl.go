package playerCoinRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	playerCoinException "github.com/Montheankul-K/game-service/pkg/playerCoin/exception"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type playerCoinRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewPlayerCoinRepositoryImpl(db databases.Database, logger echo.Logger) IPlayerCoinRepository {
	return &playerCoinRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *playerCoinRepositoryImpl) CoinAdding(tx *gorm.DB, playerCoinEntity *entities.PlayerCoin) (*entities.PlayerCoin, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	playerCoin := new(entities.PlayerCoin)
	if err := conn.Create(playerCoinEntity).Scan(playerCoin).Error; err != nil {
		r.logger.Errorf("failed to adding player coin: %s", err.Error())
		return nil, &playerCoinException.CoinAdding{}
	}
	return playerCoin, nil
}

func (r *playerCoinRepositoryImpl) Showing(playerID string) (*playerCoinModel.PlayerCoinShowing, error) {
	playerCoinShowing := new(playerCoinModel.PlayerCoinShowing)
	if err := r.db.Connect().Model(
		&entities.PlayerCoin{}).Where(
		"player_id = ?", playerID).Select(
		"player_id, sum(amount) as coin").Group(
		"player_id").Scan(playerCoinShowing).Error; err != nil {
		r.logger.Errorf("failed to showing player coin: %s", err.Error())
		return nil, &playerCoinException.PlayerCoinShowing{}
	}
	return playerCoinShowing, nil
}
