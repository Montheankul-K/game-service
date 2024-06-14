package itemShopRepository

import (
	"github.com/Montheankul-K/game-service/entities"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
)

type IItemShopRepository interface {
	Listing(itemFilter *itemShopModel.ItemFilter) ([]*entities.Item, error)
	Counting(itemFilter *itemShopModel.ItemFilter) (int64, error)
	FindByID(itemID uint64) (*entities.Item, error)
}
