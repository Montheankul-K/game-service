package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Montheankul-K/game-service/config"
	"github.com/Montheankul-K/game-service/databases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type EchoServer struct {
	app *echo.Echo
	db  databases.Database
	cfg *config.Config
}

var (
	once   sync.Once
	server *EchoServer
)

func NewEchoServer(cfg *config.Config, db databases.Database) *EchoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	once.Do(func() {
		server = &EchoServer{
			app: echoApp,
			db:  db,
			cfg: cfg,
		}
	})
	return server
}

func (s *EchoServer) Start() {
	recoverMiddleware := getRecoverMiddleware()
	loggerMiddleware := getLoggerMiddleware()
	corsMiddleware := getCORSMiddleware(s.cfg.Server.AllowOrigins)
	bodyLimitMiddleware := getBodyLimitMiddleware(s.cfg.Server.BodyLimit)
	timeoutMiddleware := getTimeoutMiddleware(s.cfg.Server.Timeout)

	s.app.Use(recoverMiddleware)
	s.app.Use(loggerMiddleware)
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	s.app.Use(timeoutMiddleware)
	authMiddleware := s.getAuthorizingMiddleware()

	s.app.GET("/v1/health", s.healthCheck)
	s.initOAuth2Router()
	s.initItemShopRouter(authMiddleware)
	s.initItemManagingRouter(authMiddleware)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go s.gracefullyShutdown(ch)
	s.listenAndServe()
	s.initPlayerCoinRouter(authMiddleware)
	s.initInventoryRouter(authMiddleware)
}

func (s *EchoServer) listenAndServe() {
	url := fmt.Sprintf(":%d", s.cfg.Server.Port)
	if err := s.app.Start(url); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.app.Logger.Fatalf("error: %s", err.Error())
	}
}

func (s *EchoServer) gracefullyShutdown(ch chan os.Signal) {
	<-ch
	s.app.Logger.Info("shutting down server")

	if err := s.app.Shutdown(context.Background()); err != nil {
		s.app.Logger.Fatalf("error: %s", err.Error())
	}
}

func (s *EchoServer) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}

func getRecoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

func getLoggerMiddleware() echo.MiddlewareFunc {
	return middleware.Logger()
}

func getTimeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "request timeout",
		Timeout:      timeout * time.Second,
	})
}

func getCORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	})
}

func getBodyLimitMiddleware(bodyLimit string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimit)
}
