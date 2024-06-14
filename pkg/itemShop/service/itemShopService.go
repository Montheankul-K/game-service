package itemShopService

import itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"

type IItemShopService interface {
	Listing(itemFilter *itemShopModel.ItemFilter) (*itemShopModel.ItemResult, error)
}
