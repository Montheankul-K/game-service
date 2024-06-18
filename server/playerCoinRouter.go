package server

import (
	playerCoinController "github.com/Montheankul-K/game-service/pkg/playerCoin/controller"
	playerCoinRepository "github.com/Montheankul-K/game-service/pkg/playerCoin/repository"
	playerCoinService "github.com/Montheankul-K/game-service/pkg/playerCoin/service"
)

func (s *EchoServer) initPlayerCoinRouter(authMiddleware *authorizingMiddleware) {
	router := s.app.Group("/v1/player-coin")
	repository := playerCoinRepository.NewPlayerCoinRepositoryImpl(s.db, s.app.Logger)
	service := playerCoinService.NewPlayerCoinServiceImpl(repository)
	controller := playerCoinController.NewPlayerCoinControllerImpl(service)

	router.POST("", controller.CoinAdding, authMiddleware.PlayerAuthorizing)
	router.GET("", controller.Showing, authMiddleware.PlayerAuthorizing)
}
