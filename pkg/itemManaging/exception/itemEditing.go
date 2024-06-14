package itemManagingException

import "fmt"

type ItemEditing struct {
	ItemID uint64
}

func (e *ItemEditing) Error() string {
	return fmt.Sprintf("edit item id: %d failed", e.ItemID)
}
