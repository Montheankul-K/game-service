package itemManagingService

import (
	itemManagingModel "github.com/Montheankul-K/game-service/pkg/itemManaging/model"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
)

type IItemManagingService interface {
	Creating(itemCreatingReq *itemManagingModel.ItemCreatingReq) (*itemShopModel.Item, error)
	Editing(itemID uint64, itemEditingReq *itemManagingModel.ItemEditingReq) (*itemShopModel.Item, error)
	Archiving(itemID uint64) error
}
