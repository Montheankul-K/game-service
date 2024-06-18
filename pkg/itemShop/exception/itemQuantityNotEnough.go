package itemShopException

import "fmt"

type ItemQuantityNotEnough struct {
	ItemID uint64
}

func (e *ItemQuantityNotEnough) Error() string {
	return fmt.Sprintf("item id:%d is not enough", e.ItemID)
}
