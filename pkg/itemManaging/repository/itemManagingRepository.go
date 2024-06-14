package itemManagingRepository

import (
	"github.com/Montheankul-K/game-service/entities"
	itemManagingModel "github.com/Montheankul-K/game-service/pkg/itemManaging/model"
)

type IItemManagingRepository interface {
	Creating(itemEntity *entities.Item) (*entities.Item, error)
	Editing(itemID uint64, itemEditingReq *itemManagingModel.ItemEditingReq) (uint64, error)
	Archiving(itemID uint64) error
}
