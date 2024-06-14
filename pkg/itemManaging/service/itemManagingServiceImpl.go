package itemManagingService

import (
	"github.com/Montheankul-K/game-service/entities"
	itemManagingModel "github.com/Montheankul-K/game-service/pkg/itemManaging/model"
	itemManagingRepository "github.com/Montheankul-K/game-service/pkg/itemManaging/repository"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
)

type itemManagingServiceImpl struct {
	itemManagingRepository itemManagingRepository.IItemManagingRepository
	itemShopRepository     itemShopRepository.IItemShopRepository
}

func NewItemManagingServiceImpl(itemManagingRepository itemManagingRepository.IItemManagingRepository, itemShopRepository itemShopRepository.IItemShopRepository) IItemManagingService {
	return &itemManagingServiceImpl{
		itemManagingRepository: itemManagingRepository,
		itemShopRepository:     itemShopRepository,
	}
}

func (s *itemManagingServiceImpl) Creating(itemCreatingReq *itemManagingModel.ItemCreatingReq) (*itemShopModel.Item, error) {
	itemEntity := &entities.Item{
		Name:        itemCreatingReq.Name,
		Description: itemCreatingReq.Description,
		Picture:     itemCreatingReq.Picture,
		Price:       itemCreatingReq.Price,
	}

	itemEntityResult, err := s.itemManagingRepository.Creating(itemEntity)
	if err != nil {
		return nil, err
	}
	return itemEntityResult.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) Editing(itemID uint64, itemEditingReq *itemManagingModel.ItemEditingReq) (*itemShopModel.Item, error) {
	_, err := s.itemManagingRepository.Editing(itemID, itemEditingReq)
	if err != nil {
		return nil, err
	}

	itemEntityResult, err := s.itemShopRepository.FindByID(itemID)
	if err != nil {
		return nil, err
	}

	return itemEntityResult.ToItemModel(), nil
}

func (s *itemManagingServiceImpl) Archiving(itemID uint64) error {
	return s.itemManagingRepository.Archiving(itemID)
}
