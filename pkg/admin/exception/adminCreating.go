package adminException

import "fmt"

type AdminCreating struct {
	AdminID string
}

func (e *AdminCreating) Error() string {
	return fmt.Sprintf("failed to create admin id: %s", e.AdminID)
}
