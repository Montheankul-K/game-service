package playerRepository

import "github.com/Montheankul-K/game-service/entities"

type IPlayerRepository interface {
	Creating(playerEntity *entities.Player) (*entities.Player, error)
	FindByID(playerID string) (*entities.Player, error)
}
