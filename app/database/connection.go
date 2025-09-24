package database

import (
	"IMP/app/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDbConnection() *gorm.DB {
	if db != nil {
		return db
	}

	dsn := "host=db user=postgres password=postgres dbname=imp port=5432 sslmode=disable"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = connection // Store the connection in the package-level variable
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		logger.Logger.Error("db is not initiated yet", nil)
	}
	return db
}
