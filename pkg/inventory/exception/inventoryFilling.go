package inventoryException

import "fmt"

type InventoryFilling struct {
	PlayerID string
	ItemID   uint64
}

func (e *InventoryFilling) Error() string {
	return fmt.Sprintf("failed to filling inventory of player id: %s with item id: %d", e.PlayerID, e.ItemID)
}
