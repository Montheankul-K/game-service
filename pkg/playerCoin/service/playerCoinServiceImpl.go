package playerCoinService

import (
	"github.com/Montheankul-K/game-service/entities"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	playerCoinRepository "github.com/Montheankul-K/game-service/pkg/playerCoin/repository"
)

type playerCoinServiceImpl struct {
	playerCoinRepository playerCoinRepository.IPlayerCoinRepository
}

func NewPlayerCoinServiceImpl(playerCoinRepository playerCoinRepository.IPlayerCoinRepository) IPlayerCoinService {
	return &playerCoinServiceImpl{
		playerCoinRepository: playerCoinRepository,
	}
}

func (s *playerCoinServiceImpl) CoinAdding(coinAddingReq *playerCoinModel.CoinAddingReq) (*playerCoinModel.PlayerCoin, error) {
	playerCoinEntity := &entities.PlayerCoin{
		PlayerID: coinAddingReq.PlayerID,
		Amount:   coinAddingReq.Amount,
	}
	playerCoinEntityResult, err := s.playerCoinRepository.CoinAdding(nil, playerCoinEntity)
	if err != nil {
		return nil, err
	}
	playerCoinEntityResult.PlayerID = coinAddingReq.PlayerID
	return playerCoinEntityResult.ToPlayerCoinModel(), nil
}

func (s *playerCoinServiceImpl) Showing(playerID string) *playerCoinModel.PlayerCoinShowing {
	playerCoinShowing, err := s.playerCoinRepository.Showing(playerID)
	if err != nil {
		return &playerCoinModel.PlayerCoinShowing{
			PlayerID: playerID,
			Coin:     0,
		}
	}
	return playerCoinShowing
}
