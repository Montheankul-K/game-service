package itemShopRepository

import (
	"github.com/Montheankul-K/game-service/entities"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"gorm.io/gorm"
)

type IItemShopRepository interface {
	TransactionBegin() *gorm.DB
	TransactionRollback(tx *gorm.DB) error
	TransactionCommit(tx *gorm.DB) error
	Listing(itemFilter *itemShopModel.ItemFilter) ([]*entities.Item, error)
	Counting(itemFilter *itemShopModel.ItemFilter) (int64, error)
	FindByID(itemID uint64) (*entities.Item, error)
	FindByIDList(itemIDs []uint64) ([]*entities.Item, error)
	PurchaseHistoryRecording(tx *gorm.DB, purchasingEntity *entities.PurchaseHistory) (*entities.PurchaseHistory, error)
}
