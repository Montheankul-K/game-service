package playerRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	playerException "github.com/Montheankul-K/game-service/pkg/player/exception"
	"github.com/labstack/echo/v4"
)

type playerRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewPlayerRepositoryImpl(db databases.Database, logger echo.Logger) IPlayerRepository {
	return &playerRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *playerRepositoryImpl) Creating(playerEntity *entities.Player) (*entities.Player, error) {
	player := new(entities.Player)
	if err := r.db.Connect().Create(playerEntity).Scan(player).Error; err != nil {
		r.logger.Errorf("failed to creating player: %s", err)
		return nil, &playerException.PlayerCreating{
			PlayerID: playerEntity.ID,
		}
	}
	return player, nil
}

func (r *playerRepositoryImpl) FindByID(playerID string) (*entities.Player, error) {
	player := new(entities.Player)
	if err := r.db.Connect().Where("id = ?", playerID).First(player).Error; err != nil {
		r.logger.Errorf("failed to finding player id: %s", err)
		return nil, &playerException.PlayerNotFound{
			PlayerID: playerID,
		}
	}
	return player, nil
}
