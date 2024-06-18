package inventoryService

import inventoryModel "github.com/Montheankul-K/game-service/pkg/inventory/model"

type IInventoryService interface {
	Listing(playerID string) ([]*inventoryModel.Inventory, error)
}
