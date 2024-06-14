package oauth2Service

import (
	adminModel "github.com/Montheankul-K/game-service/pkg/admin/model"
	playerModel "github.com/Montheankul-K/game-service/pkg/player/model"
)

type IOAuth2Service interface {
	PlayerAccountCreating(playerCreatingReq *playerModel.PlayerCreatingReq) error
	AdminAccountCreating(adminCreatingReq *adminModel.AdminCreatingReq) error
	IsPlayer(playerID string) bool
	IsAdmin(adminID string) bool
}
