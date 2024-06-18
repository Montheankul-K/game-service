package itemShopService

import (
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
)

type IItemShopService interface {
	Listing(itemFilter *itemShopModel.ItemFilter) (*itemShopModel.ItemResult, error)
	Buying(buyingReq *itemShopModel.BuyingReq) (*playerCoinModel.PlayerCoin, error)
	Selling(sellingReq *itemShopModel.SellingReq) (*playerCoinModel.PlayerCoin, error)
}
