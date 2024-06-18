package itemShopException

type HistoryOfPurchaseRecording struct{}

func (e *HistoryOfPurchaseRecording) Error() string {
	return "failed to recording history of purchase"
}
