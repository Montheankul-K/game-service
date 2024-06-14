package server

import (
	itemManagingController "github.com/Montheankul-K/game-service/pkg/itemManaging/controller"
	itemManagingRepository "github.com/Montheankul-K/game-service/pkg/itemManaging/repository"
	itemManagingService "github.com/Montheankul-K/game-service/pkg/itemManaging/service"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
)

func (s *EchoServer) initItemManagingRouter(authMiddleware *authorizingMiddleware) {
	router := s.app.Group("/v1/item-managing")
	itemManagingRepo := itemManagingRepository.NewItemManagingRepositoryImpl(s.db, s.app.Logger)
	itemShopRepo := itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)
	service := itemManagingService.NewItemManagingServiceImpl(itemManagingRepo, itemShopRepo)
	controller := itemManagingController.NewItemManagingController(service)

	router.POST("", controller.Creating, authMiddleware.AdminAuthorizing)
	router.PATCH("/:itemID", controller.Editing, authMiddleware.AdminAuthorizing)
	router.DELETE("/:itemID", controller.Archiving, authMiddleware.AdminAuthorizing)
}
