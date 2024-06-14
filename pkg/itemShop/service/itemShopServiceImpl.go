package itemShopService

import (
	"github.com/Montheankul-K/game-service/entities"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"github.com/Montheankul-K/game-service/pkg/itemShop/repository"
)

type itemShopServiceImpl struct {
	itemShopRepository itemShopRepository.IItemShopRepository
}

func NewItemShopServiceImpl(itemShopRepository itemShopRepository.IItemShopRepository) IItemShopService {
	return &itemShopServiceImpl{
		itemShopRepository: itemShopRepository,
	}
}

func (s *itemShopServiceImpl) Listing(itemFilter *itemShopModel.ItemFilter) (*itemShopModel.ItemResult, error) {
	itemList, err := s.itemShopRepository.Listing(itemFilter)
	if err != nil {
		return nil, err
	}

	itemCounting, err := s.itemShopRepository.Counting(itemFilter)
	if err != nil {
		return nil, err
	}
	totalPage := s.totalPageCalculation(itemCounting, itemFilter.Size)

	return s.toItemResultResponse(itemList, itemFilter.Page, totalPage), nil
}

func (s *itemShopServiceImpl) totalPageCalculation(totalItem int64, size int64) int64 {
	totalPage := totalItem / size
	if totalItem%size != 0 {
		totalPage++
	}
	return totalPage
}

func (s *itemShopServiceImpl) toItemResultResponse(itemEntityList []*entities.Item, page, totalPage int64) *itemShopModel.ItemResult {
	itemModelList := make([]*itemShopModel.Item, len(itemEntityList))
	for _, item := range itemEntityList {
		itemModelList = append(itemModelList, item.ToItemModel())
	}
	return &itemShopModel.ItemResult{
		Items: itemModelList,
		Paginate: itemShopModel.PaginateResult{
			Page:      page,
			TotalPage: totalPage,
		},
	}
}
