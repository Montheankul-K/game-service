package server

import (
	inventoryRepository "github.com/Montheankul-K/game-service/pkg/inventory/repository"
	itemShopController "github.com/Montheankul-K/game-service/pkg/itemShop/controller"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
	itemShopService "github.com/Montheankul-K/game-service/pkg/itemShop/service"
	playerCoinRepository "github.com/Montheankul-K/game-service/pkg/playerCoin/repository"
)

func (s *EchoServer) initItemShopRouter(authMiddleware *authorizingMiddleware) {
	router := s.app.Group("/v1/item-shop")
	itemShopRepo := itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)
	playerCoinRepo := playerCoinRepository.NewPlayerCoinRepositoryImpl(s.db, s.app.Logger)
	inventoryRepo := inventoryRepository.NewInventoryRepositoryImpl(s.db, s.app.Logger)
	service := itemShopService.NewItemShopServiceImpl(itemShopRepo, playerCoinRepo, inventoryRepo, s.app.Logger)
	controller := itemShopController.NewItemShopController(service)

	router.GET("/listing", controller.Listing, authMiddleware.PlayerAuthorizing)
	router.POST("/buying", controller.Buying, authMiddleware.PlayerAuthorizing)
	router.POST("/selling", controller.Selling, authMiddleware.PlayerAuthorizing)
}
