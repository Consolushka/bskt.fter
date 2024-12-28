package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connect() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=imp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}
