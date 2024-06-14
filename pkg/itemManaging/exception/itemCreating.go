package itemManagingException

type ItemCreating struct{}

func (e *ItemCreating) Error() string {
	return "create item failed"
}
