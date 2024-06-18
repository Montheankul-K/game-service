package inventoryModel

import itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"

type Inventory struct {
	Item     *itemShopModel.Item `json:"item"`
	Quantity uint                `json:"quantity"`
}

type ItemQuantityCounting struct {
	ItemID   uint64
	Quantity uint
}
