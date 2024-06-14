package databases

import (
	"fmt"
	"github.com/Montheankul-K/game-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"sync"
)

type postgresDatabase struct {
	*gorm.DB
}

var (
	postgresDBInstance *postgresDatabase
	once               sync.Once
)

func NewPostgresDatabase(cfg *config.Database) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s search_path=%s",
			cfg.Host,
			cfg.Username,
			cfg.Password,
			cfg.DBName,
			cfg.Port,
			cfg.SSLMode,
			cfg.Schema,
		)

		conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		log.Printf("connected to databases: %s", cfg.DBName)
		postgresDBInstance = &postgresDatabase{
			DB: conn,
		}
	})
	return postgresDBInstance
}

func (db *postgresDatabase) Connect() *gorm.DB {
	return postgresDBInstance.DB
}
