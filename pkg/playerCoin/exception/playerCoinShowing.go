package playerCoinException

type PlayerCoinShowing struct{}

func (e *PlayerCoinShowing) Error() string {
	return "failed to showing coin"
}
