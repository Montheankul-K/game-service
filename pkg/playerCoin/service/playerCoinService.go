package playerCoinService

import playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"

type IPlayerCoinService interface {
	CoinAdding(coinAddingReq *playerCoinModel.CoinAddingReq) (*playerCoinModel.PlayerCoin, error)
	Showing(playerID string) *playerCoinModel.PlayerCoinShowing
}
