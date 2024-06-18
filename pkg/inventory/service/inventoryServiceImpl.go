package inventoryService

import (
	"github.com/Montheankul-K/game-service/entities"
	inventoryModel "github.com/Montheankul-K/game-service/pkg/inventory/model"
	inventoryRepository "github.com/Montheankul-K/game-service/pkg/inventory/repository"
	itemShopRepository "github.com/Montheankul-K/game-service/pkg/itemShop/repository"
)

type inventoryServiceImpl struct {
	inventoryRepository inventoryRepository.IInventoryRepository
	itemShopRepository  itemShopRepository.IItemShopRepository
}

func NewInventoryServiceImpl(inventoryRepository inventoryRepository.IInventoryRepository, itemShopRepository itemShopRepository.IItemShopRepository) IInventoryService {
	return &inventoryServiceImpl{
		inventoryRepository: inventoryRepository,
		itemShopRepository:  itemShopRepository,
	}
}

func (s *inventoryServiceImpl) Listing(playerID string) ([]*inventoryModel.Inventory, error) {
	inventoryEntities, err := s.inventoryRepository.Listing(playerID)
	if err != nil {
		return nil, err
	}

	uniqueItemWithQuantityCounterList := s.getUniqueItemWithQuantityCounterList(inventoryEntities)
	return s.buildInventoryListingResult(uniqueItemWithQuantityCounterList), nil
}

func (s *inventoryServiceImpl) getUniqueItemWithQuantityCounterList(inventoryEntities []*entities.Inventory) []inventoryModel.ItemQuantityCounting {
	itemQuantityCounterList := make([]inventoryModel.ItemQuantityCounting, 0)
	itemMapWithQuantity := make(map[uint64]uint)
	for _, inventory := range inventoryEntities {
		itemMapWithQuantity[inventory.ItemID]++
	}

	for itemID, quantity := range itemMapWithQuantity {
		itemQuantityCounterList = append(itemQuantityCounterList, inventoryModel.ItemQuantityCounting{
			ItemID:   itemID,
			Quantity: quantity,
		})
	}
	return itemQuantityCounterList
}

func (s *inventoryServiceImpl) buildInventoryListingResult(uniqueItemWithQuantityCounterList []inventoryModel.ItemQuantityCounting) []*inventoryModel.Inventory {
	uniqueItemIDList := s.getItemID(uniqueItemWithQuantityCounterList)
	itemEntities, err := s.itemShopRepository.FindByIDList(uniqueItemIDList)
	if err != nil {
		return make([]*inventoryModel.Inventory, 0)
	}

	results := make([]*inventoryModel.Inventory, 0)
	itemMapWithQuantity := s.getItemMapWithQuantity(uniqueItemWithQuantityCounterList)
	for _, itemEntity := range itemEntities {
		results = append(results, &inventoryModel.Inventory{
			Item:     itemEntity.ToItemModel(),
			Quantity: itemMapWithQuantity[itemEntity.ID],
		})
	}
	return results
}

func (s *inventoryServiceImpl) getItemID(uniqueItemWithQuantityCounterList []inventoryModel.ItemQuantityCounting) []uint64 {
	uniqueItemIDList := make([]uint64, len(uniqueItemWithQuantityCounterList))
	for _, inventory := range uniqueItemWithQuantityCounterList {
		uniqueItemIDList = append(uniqueItemIDList, inventory.ItemID)
	}
	return uniqueItemIDList
}

func (s *inventoryServiceImpl) getItemMapWithQuantity(uniqueItemWithQuantityCounterList []inventoryModel.ItemQuantityCounting) map[uint64]uint {
	itemMapWithQuantity := make(map[uint64]uint)
	for _, inventory := range uniqueItemWithQuantityCounterList {
		itemMapWithQuantity[inventory.ItemID] = inventory.Quantity
	}
	return itemMapWithQuantity
}
