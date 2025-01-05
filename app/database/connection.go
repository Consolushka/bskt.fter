package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
	fmt.Println("Successfully connected to database")
	return db
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Panic("Database connection is nil")
		return OpenDbConnection()
	}
	return db
}
