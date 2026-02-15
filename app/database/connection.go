package database

import (
	"IMP/app/pkg/logger"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDbConnection() *gorm.DB {
	if db != nil {
		return db
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("failed to connect database", map[string]interface{}{
			"host":   host,
			"dbName": dbName,
			"port":   port,
			"error":  err,
		})
		panic("failed to connect database")
	}

	db = connection // Store the connection in the package-level variable
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		logger.Error("db is not initiated yet", nil)
	}
	return db
}
