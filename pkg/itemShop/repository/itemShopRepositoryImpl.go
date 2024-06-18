package itemShopRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	itemShopException "github.com/Montheankul-K/game-service/pkg/itemShop/exception"
	itemShopModel "github.com/Montheankul-K/game-service/pkg/itemShop/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type itemShopRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemShopRepositoryImpl(db databases.Database, logger echo.Logger) IItemShopRepository {
	return &itemShopRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *itemShopRepositoryImpl) TransactionBegin() *gorm.DB {
	tx := r.db.Connect()
	return tx.Begin()
}

func (r *itemShopRepositoryImpl) TransactionRollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *itemShopRepositoryImpl) TransactionCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *itemShopRepositoryImpl) Listing(itemFilter *itemShopModel.ItemFilter) ([]*entities.Item, error) {
	itemList := make([]*entities.Item, 0)
	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false)
	if itemFilter.Name != "" {
		query = query.Where("name LIKE ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description LIKE ?", "%"+itemFilter.Description+"%")
	}

	offset := int((itemFilter.Page - 1) * itemFilter.Size)
	if err := query.Offset(offset).Limit(int(itemFilter.Size)).Find(&itemList).Order("id asc").Error; err != nil {
		r.logger.Errorf("failed to list items: %s", err.Error())
		return nil, &itemShopException.ItemListing{}
	}
	return itemList, nil
}

func (r *itemShopRepositoryImpl) Counting(itemFilter *itemShopModel.ItemFilter) (int64, error) {
	query := r.db.Connect().Model(&entities.Item{}).Where("is_archive = ?", false)
	if itemFilter.Name != "" {
		query = query.Where("name LIKE ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description LIKE ?", "%"+itemFilter.Description+"%")
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		r.logger.Errorf("failed to counting item: %s", err.Error())
		return -1, &itemShopException.ItemCounting{}
	}
	return count, nil
}

func (r *itemShopRepositoryImpl) FindByID(itemID uint64) (*entities.Item, error) {
	item := new(entities.Item)
	if err := r.db.Connect().First(item, itemID).Error; err != nil {
		r.logger.Errorf("failed to find item: %s", err.Error())
		return nil, &itemShopException.ItemNotFound{
			ItemID: itemID,
		}
	}
	return item, nil
}

func (r *itemShopRepositoryImpl) FindByIDList(itemIDs []uint64) ([]*entities.Item, error) {
	item := make([]*entities.Item, len(itemIDs))
	if err := r.db.Connect().Model(&entities.Item{}).Where("id IN (?)", itemIDs).Find(&item).Error; err != nil {
		r.logger.Errorf("failed to find items by ID list: %s", err.Error())
		return nil, &itemShopException.ItemListing{}
	}
	return item, nil
}

func (r *itemShopRepositoryImpl) PurchaseHistoryRecording(tx *gorm.DB, purchasingEntity *entities.PurchaseHistory) (*entities.PurchaseHistory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	insertedPurchasing := new(entities.PurchaseHistory)
	if err := conn.Create(purchasingEntity).Error; err != nil {
		r.logger.Errorf("failed to record purchase history: %s", err.Error())
		return nil, &itemShopException.HistoryOfPurchaseRecording{}
	}
	return insertedPurchasing, nil
}
