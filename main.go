package main

import (
	"github.com/Montheankul-K/game-service/config"
	"github.com/Montheankul-K/game-service/databases"
	"github.com/Montheankul-K/game-service/server"
)

func main() {
	cfg := config.GetConfig()
	db := databases.NewPostgresDatabase(cfg.Database)
	srv := server.NewEchoServer(cfg, db)

	srv.Start()
}
