package itemManagingException

import "fmt"

type ArchiveItem struct {
	ItemID uint64
}

func (e *ArchiveItem) Error() string {
	return fmt.Sprintf("archive item id: %d failed", e.ItemID)
}
