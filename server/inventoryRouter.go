package server

import (
	inventorycontroller "github.com/Montheankul-K/game-service/pkg/inventory/controller"
	inventoryRepository "github.com/Montheankul-K/game-service/pkg/inventory/repository"
	inventoryService "github.com/Montheankul-K/game-service/pkg/inventory/service"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
)

func (s *EchoServer) initInventoryRouter(authMiddleware *authorizingMiddleware) {
	router := s.app.Group("/v1/inventory")
	inventoryRepo := inventoryRepository.NewInventoryRepositoryImpl(s.db, s.app.Logger)
	itemShopRepo := itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)
	service := inventoryService.NewInventoryServiceImpl(inventoryRepo, itemShopRepo)
	controller := inventorycontroller.NewInventoryControllerImpl(service, s.app.Logger)

	router.GET("", controller.Listing, authMiddleware.PlayerAuthorizing)
}
