package playerCoinException

type CoinAdding struct{}

func (e *CoinAdding) Error() string {
	return "failed to adding coin"
}
