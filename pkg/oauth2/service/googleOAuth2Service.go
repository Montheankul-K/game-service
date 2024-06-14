package oauth2Service

import (
	"github.com/Montheankul-K/game-service/entities"
	adminModel "github.com/Montheankul-K/game-service/pkg/admin/model"
	adminRepository "github.com/Montheankul-K/game-service/pkg/admin/repository"
	playerModel "github.com/Montheankul-K/game-service/pkg/player/model"
	playerRepository "github.com/Montheankul-K/game-service/pkg/player/repository"
)

type googleOAuth2Service struct {
	playerRepository playerRepository.IPlayerRepository
	adminRepository  adminRepository.IAdminRepository
}

func NewGoogleOAuth2Service(playerRepository playerRepository.IPlayerRepository, adminRepository adminRepository.IAdminRepository) IOAuth2Service {
	return &googleOAuth2Service{
		playerRepository: playerRepository,
		adminRepository:  adminRepository,
	}
}

func (s *googleOAuth2Service) PlayerAccountCreating(playerCreatingReq *playerModel.PlayerCreatingReq) error {
	if !s.IsPlayer(playerCreatingReq.ID) {
		playerEntity := &entities.Player{
			ID:     playerCreatingReq.ID,
			Name:   playerCreatingReq.Name,
			Email:  playerCreatingReq.Email,
			Avatar: playerCreatingReq.Avatar,
		}

		if _, err := s.playerRepository.Creating(playerEntity); err != nil {
			return err
		}
	}
	return nil
}

func (s *googleOAuth2Service) AdminAccountCreating(adminCreatingReq *adminModel.AdminCreatingReq) error {
	if !s.IsAdmin(adminCreatingReq.ID) {
		adminEntity := &entities.Admin{
			ID:     adminCreatingReq.ID,
			Name:   adminCreatingReq.Name,
			Email:  adminCreatingReq.Email,
			Avatar: adminCreatingReq.Avatar,
		}

		if _, err := s.adminRepository.Creating(adminEntity); err != nil {
			return err
		}
	}
	return nil
}

func (s *googleOAuth2Service) IsPlayer(playerID string) bool {
	player, err := s.playerRepository.FindByID(playerID)
	if err != nil {
		return false
	}
	return player != nil
}

func (s *googleOAuth2Service) IsAdmin(adminID string) bool {
	admin, err := s.adminRepository.FindByID(adminID)
	if err != nil {
		return false
	}
	return admin != nil
}
