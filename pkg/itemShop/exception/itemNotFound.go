package itemShopException

import "fmt"

type ItemNotFound struct {
	ItemID uint64
}

func (e *ItemNotFound) Error() string {
	return fmt.Sprintf("item id: %d not found", e.ItemID)
}
