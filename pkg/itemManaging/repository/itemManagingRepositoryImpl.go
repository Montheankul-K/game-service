package itemManagingRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	itemManagingException "github.com/Montheankul-K/game-service/pkg/itemManaging/exception"
	itemManagingModel "github.com/Montheankul-K/game-service/pkg/itemManaging/model"
	"github.com/labstack/echo/v4"
)

type itemManagingRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemManagingRepositoryImpl(db databases.Database, logger echo.Logger) IItemManagingRepository {
	return &itemManagingRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *itemManagingRepositoryImpl) Creating(itemEntity *entities.Item) (*entities.Item, error) {
	item := new(entities.Item)
	if err := r.db.Connect().Create(itemEntity).Scan(item).Error; err != nil {
		r.logger.Errorf("failed to creating item: %s", err)
		return nil, &itemManagingException.ItemCreating{}
	}
	return item, nil
}

func (r *itemManagingRepositoryImpl) Editing(itemID uint64, itemEditingReq *itemManagingModel.ItemEditingReq) (uint64, error) {
	if err := r.db.Connect().Model(&entities.Item{}).Where("id =?", itemID).Updates(itemEditingReq).Error; err != nil {
		r.logger.Errorf("failed to editing item: %s", err)
		return 0, &itemManagingException.ItemEditing{
			ItemID: itemID,
		}
	}
	return itemID, nil
}

func (r *itemManagingRepositoryImpl) Archiving(itemID uint64) error {
	if err := r.db.Connect().Table("items").Where("id =?", itemID).Update("is_archive", true).Error; err != nil {
		r.logger.Errorf("failed to archiving item: %s", err)
		return &itemManagingException.ArchiveItem{
			ItemID: itemID,
		}
	}
	return nil
}
