package inventoryRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	inventoryException "github.com/Montheankul-K/game-service/pkg/inventory/exception"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type inventoryRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewInventoryRepositoryImpl(db databases.Database, logger echo.Logger) IInventoryRepository {
	return &inventoryRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *inventoryRepositoryImpl) Filling(tx *gorm.DB, playerID string, itemID uint64, qty int) ([]*entities.Inventory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	inventoryEntities := make([]*entities.Inventory, 0)
	for range qty {
		inventoryEntities = append(inventoryEntities, &entities.Inventory{
			PlayerID: playerID,
			ItemID:   itemID,
		})
	}

	if err := conn.CreateInBatches(inventoryEntities, len(inventoryEntities)).Error; err != nil {
		r.logger.Errorf("failed to filling inventory: %s", err.Error())
		return nil, &inventoryException.InventoryFilling{
			PlayerID: playerID,
			ItemID:   itemID,
		}
	}
	return inventoryEntities, nil
}

func (r *inventoryRepositoryImpl) Removing(tx *gorm.DB, playerID string, itemID uint64, limit int) error {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	inventoryEntities, err := r.findPlayerItemInInventoryByID(playerID, itemID, limit)
	if err != nil {
		return err
	}

	for _, inventory := range inventoryEntities {
		inventory.IsDeleted = true
		if err = conn.Model(&entities.Inventory{}).Where("id = ?", inventory.ID).Updates(inventory).Error; err != nil {
			r.logger.Errorf("failed to removing item in inventory: %s", err.Error())
			return &inventoryException.PlayerItemRemoving{
				ItemID: itemID,
			}
		}
	}
	return nil
}

func (r *inventoryRepositoryImpl) PlayerItemCounting(playerID string, itemID uint64) int64 {
	var count int64
	if err := r.db.Connect().Model(&entities.Inventory{}).Where(
		"player_id = ? AND item_id = ? AND is_deleted = ?", playerID, itemID, false).Count(&count).Error; err != nil {
		r.logger.Errorf("failed to counting item in inventory: %s", err.Error())
		return -1
	}
	return count
}

func (r *inventoryRepositoryImpl) Listing(playerID string) ([]*entities.Inventory, error) {
	inventoryEntities := make([]*entities.Inventory, 0)
	if err := r.db.Connect().Where("player_id = ? AND is_deleted = ?", playerID, false).Find(&inventoryEntities).Error; err != nil {
		r.logger.Errorf("failed to listing player inventory: %s", err.Error())
		return nil, &inventoryException.PlayerItemFinding{
			PlayerID: playerID,
		}
	}
	return inventoryEntities, nil
}

func (r *inventoryRepositoryImpl) findPlayerItemInInventoryByID(playerID string, itemID uint64, limit int) ([]*entities.Inventory, error) {
	inventoryEntities := make([]*entities.Inventory, 0)
	if err := r.db.Connect().Where(
		"player_id = ? AND item_id = ? AND is_deleted = ?", playerID, itemID, false).Limit(
		limit).Find(&inventoryEntities).Error; err != nil {
		r.logger.Errorf("failed to find player item in inventory: %s", err.Error())
		return nil, &inventoryException.PlayerItemRemoving{
			ItemID: itemID,
		}
	}
	return inventoryEntities, nil
}
