package inventoryRepository

import (
	"github.com/Montheankul-K/game-service/entities"
	"gorm.io/gorm"
)

type IInventoryRepository interface {
	Filling(tx *gorm.DB, playerID string, itemID uint64, qty int) ([]*entities.Inventory, error)
	Removing(tx *gorm.DB, playerID string, itemID uint64, limit int) error
	PlayerItemCounting(playerID string, itemID uint64) int64
	Listing(playerID string) ([]*entities.Inventory, error)
}
