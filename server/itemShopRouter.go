package server

import (
	itemShopController "github.com/Montheankul-K/game-service/pkg/itemShop/controller"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
	itemShopService "github.com/Montheankul-K/game-service/pkg/itemShop/service"
)

func (s *EchoServer) initItemShopRouter() {
	router := s.app.Group("/v1/item-shop")
	repository := itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)
	service := itemShopService.NewItemShopServiceImpl(repository)
	controller := itemShopController.NewItemShopController(service)

	router.GET("/listing", controller.Listing)
}
