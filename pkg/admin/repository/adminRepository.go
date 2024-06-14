package adminRepository

import "github.com/Montheankul-K/game-service/entities"

type IAdminRepository interface {
	Creating(adminEntity *entities.Admin) (*entities.Admin, error)
	FindByID(adminID string) (*entities.Admin, error)
}
