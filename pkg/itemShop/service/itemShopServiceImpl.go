package itemShopService

import (
	"github.com/Montheankul-K/game-service/entities"
	inventoryRepository "github.com/Montheankul-K/game-service/pkg/inventory/repository"
	itemShopException "github.com/Montheankul-K/game-service/pkg/itemShop/exception"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"github.com/Montheankul-K/game-service/pkg/itemShop/repository"
	playerCoinModel "github.com/Montheankul-K/game-service/pkg/playerCoin/model"
	playerCoinRepository "github.com/Montheankul-K/game-service/pkg/playerCoin/repository"
	"github.com/labstack/echo/v4"
)

type itemShopServiceImpl struct {
	itemShopRepository   itemShopRepository.IItemShopRepository
	playerCoinRepository playerCoinRepository.IPlayerCoinRepository
	inventoryRepository  inventoryRepository.IInventoryRepository
	logger               echo.Logger
}

func NewItemShopServiceImpl(itemShopRepository itemShopRepository.IItemShopRepository, playerCoinRepository playerCoinRepository.IPlayerCoinRepository, inventoryRepository inventoryRepository.IInventoryRepository, logger echo.Logger) IItemShopService {
	return &itemShopServiceImpl{
		itemShopRepository:   itemShopRepository,
		playerCoinRepository: playerCoinRepository,
		inventoryRepository:  inventoryRepository,
		logger:               logger,
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

func (s *itemShopServiceImpl) Buying(buyingReq *itemShopModel.BuyingReq) (*playerCoinModel.PlayerCoin, error) {
	itemEntity, err := s.itemShopRepository.FindByID(buyingReq.ItemID)
	if err != nil {
		return nil, err
	}

	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), buyingReq.Quantity)
	if err = s.playerCoinChecking(buyingReq.PlayerID, totalPrice); err != nil {
		return nil, err
	}

	tx := s.itemShopRepository.TransactionBegin()
	purchaseRecording, err := s.itemShopRepository.PurchaseHistoryRecording(tx, &entities.PurchaseHistory{
		PlayerID:        buyingReq.PlayerID,
		ItemID:          buyingReq.ItemID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        buyingReq.Quantity,
		IsBuying:        true,
	})
	if err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("purchase history recorded: %d", purchaseRecording.ID)

	playerCoin, err := s.playerCoinRepository.CoinAdding(tx, &entities.PlayerCoin{
		PlayerID: buyingReq.PlayerID,
		Amount:   -totalPrice,
	})
	if err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("player coin deducted: %d", playerCoin.Amount)

	inventoryEntity, err := s.inventoryRepository.Filling(tx, buyingReq.PlayerID, buyingReq.ItemID, int(buyingReq.Quantity))
	if err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("inventory filled: %d", len(inventoryEntity))

	if err = s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}
	return playerCoin.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) Selling(sellingReq *itemShopModel.SellingReq) (*playerCoinModel.PlayerCoin, error) {
	itemEntity, err := s.itemShopRepository.FindByID(sellingReq.ItemID)
	if err != nil {
		return nil, err
	}

	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), sellingReq.Quantity)
	totalPrice = totalPrice / 2
	if err = s.playerItemChecking(sellingReq.PlayerID, sellingReq.ItemID, sellingReq.Quantity); err != nil {
		return nil, err
	}

	tx := s.itemShopRepository.TransactionBegin()
	purchaseRecording, err := s.itemShopRepository.PurchaseHistoryRecording(tx, &entities.PurchaseHistory{
		PlayerID:        sellingReq.PlayerID,
		ItemID:          sellingReq.ItemID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        sellingReq.Quantity,
		IsBuying:        false,
	})
	if err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("selling history recorded: %d", purchaseRecording.ID)

	playerCoin, err := s.playerCoinRepository.CoinAdding(tx, &entities.PlayerCoin{
		PlayerID: sellingReq.PlayerID,
		Amount:   totalPrice,
	})
	if err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("player coin added: %d", playerCoin.Amount)

	if err = s.inventoryRepository.Removing(tx, sellingReq.PlayerID, sellingReq.ItemID, int(sellingReq.Quantity)); err != nil {
		_ = s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("inventory item id: %d removed: %d", sellingReq.ItemID, sellingReq.Quantity)

	if err = s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}
	return playerCoin.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) totalPriceCalculation(item *itemShopModel.Item, qty uint) int64 {
	return int64(item.Price) * int64(qty)
}

func (s *itemShopServiceImpl) playerCoinChecking(playerID string, totalPrice int64) error {
	playerCoin, err := s.playerCoinRepository.Showing(playerID)
	if err != nil {
		return err
	}

	if playerCoin.Coin < totalPrice {
		s.logger.Error("player coin is not enough")
		return &itemShopException.CoinNotEnough{}
	}
	return nil
}

func (s *itemShopServiceImpl) playerItemChecking(playerID string, itemID uint64, qty uint) error {
	itemCounting := s.inventoryRepository.PlayerItemCounting(playerID, itemID)
	if int(itemCounting) < int(qty) {
		return &itemShopException.ItemQuantityNotEnough{
			ItemID: itemID,
		}
	}
	return nil
}
