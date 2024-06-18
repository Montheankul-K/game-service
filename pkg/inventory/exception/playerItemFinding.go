package inventoryException

import "fmt"

type PlayerItemFinding struct {
	PlayerID string
}

func (e *PlayerItemFinding) Error() string {
	return fmt.Sprintf("failed to find item for player id: %s", e.PlayerID)
}
