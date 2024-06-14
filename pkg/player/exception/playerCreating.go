package playerException

import "fmt"

type PlayerCreating struct {
	PlayerID string
}

func (e *PlayerCreating) Error() string {
	return fmt.Sprintf("failed to create player id: %s", e.PlayerID)
}
