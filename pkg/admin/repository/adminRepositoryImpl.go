package adminRepository

import (
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/entities"
	adminException "github.com/Montheankul-K/game-service/pkg/admin/exception"
	"github.com/labstack/echo/v4"
)

type adminRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewAdminRepositoryImpl(db databases.Database, logger echo.Logger) IAdminRepository {
	return &adminRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *adminRepositoryImpl) Creating(adminEntity *entities.Admin) (*entities.Admin, error) {
	admin := new(entities.Admin)
	if err := r.db.Connect().Create(adminEntity).Scan(admin).Error; err != nil {
		r.logger.Errorf("failed to creating player: %s", err)
		return nil, &adminException.AdminCreating{
			AdminID: adminEntity.ID,
		}
	}
	return admin, nil
}

func (r *adminRepositoryImpl) FindByID(adminID string) (*entities.Admin, error) {
	admin := new(entities.Admin)
	if err := r.db.Connect().Where("id = ?", adminID).First(admin).Error; err != nil {
		r.logger.Errorf("failed to finding player id: %s", err)
		return nil, &adminException.AdminNotFound{
			AdminID: adminID,
		}
	}
	return admin, nil
}
