package server

import (
	"github.com/Montheankul-K/game-service/config"
	adminRepository "github.com/Montheankul-K/game-service/pkg/admin/repository"
	oauth2controller "github.com/Montheankul-K/game-service/pkg/oauth2/controller"
	oauth2Service "github.com/Montheankul-K/game-service/pkg/oauth2/service"
	playerRepository "github.com/Montheankul-K/game-service/pkg/player/repository"
	"github.com/labstack/echo/v4"
)

type authorizingMiddleware struct {
	oauth2Controller oauth2controller.IOAuth2Controller
	oauth2Conf       *config.OAuth2
	logger           echo.Logger
}

func (m *authorizingMiddleware) PlayerAuthorizing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return m.oauth2Controller.PlayerAuthorizing(ctx, next)
	}
}

func (m *authorizingMiddleware) AdminAuthorizing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return m.oauth2Controller.AdminAuthorizing(ctx, next)
	}
}

func (s *EchoServer) getAuthorizingMiddleware() *authorizingMiddleware {
	playerRepo := playerRepository.NewPlayerRepositoryImpl(s.db, s.app.Logger)
	adminRepo := adminRepository.NewAdminRepositoryImpl(s.db, s.app.Logger)
	service := oauth2Service.NewGoogleOAuth2Service(playerRepo, adminRepo)
	controller := oauth2controller.NewGoogleOAuth2Controller(service, s.cfg.OAuth2, s.app.Logger)

	return &authorizingMiddleware{
		oauth2Controller: controller,
		oauth2Conf:       s.cfg.OAuth2,
		logger:           s.app.Logger,
	}
}
