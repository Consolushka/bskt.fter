package main

import (
	"IMP/app/database"
	"IMP/app/internal/service/scheduler"
	"IMP/app/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC

	logger.Init(logger.BuildLoggers())

	if err := godotenv.Load(); err != nil {
		logger.Error("Couldn't load .env file", map[string]interface{}{
			"error": err,
		})
		panic(err)
	}

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
