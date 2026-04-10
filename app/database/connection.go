package database

import (
	"IMP/app/internal/infra/config"
	"fmt"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDbConnection(cfg config.DatabaseConfig) *gorm.DB {
	if db != nil {
		return db
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		compositelogger.Error("failed to connect database", map[string]interface{}{
			"host":   cfg.Host,
			"dbName": cfg.Name,
			"port":   cfg.Port,
			"error":  err,
		})
		panic("failed to connect database")
	}

	db = connection // Store the connection in the package-level variable
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		compositelogger.Error("db is not initiated yet", nil)
	}
	return db
}
