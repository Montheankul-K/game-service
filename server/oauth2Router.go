package server

import (
	adminRepository "github.com/Montheankul-K/game-service/pkg/admin/repository"
	oauth2controller "github.com/Montheankul-K/game-service/pkg/oauth2/controller"
	oauth2Service "github.com/Montheankul-K/game-service/pkg/oauth2/service"
	playerRepository "github.com/Montheankul-K/game-service/pkg/player/repository"
)

func (s *EchoServer) initOAuth2Router() {
	router := s.app.Group("/v1/oauth2/google")
	playerRepo := playerRepository.NewPlayerRepositoryImpl(s.db, s.app.Logger)
	adminRepo := adminRepository.NewAdminRepositoryImpl(s.db, s.app.Logger)
	service := oauth2Service.NewGoogleOAuth2Service(playerRepo, adminRepo)
	controller := oauth2controller.NewGoogleOAuth2Controller(service, s.cfg.OAuth2, s.app.Logger)

	router.GET("/player/login", controller.PlayerLogin)
	router.GET("/admin/login", controller.AdminLogin)
	router.GET("/player/login/callback", controller.PlayerLoginCallback)
	router.GET("/admin/login/callback", controller.AdminLoginCallback)
	router.DELETE("/logout", controller.Logout)
}
