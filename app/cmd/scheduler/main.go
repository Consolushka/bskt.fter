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

	if err := godotenv.Load(); err != nil {
		logger.Error("Couldn't load .env file", map[string]interface{}{
			"error": err,
		})
		panic(err)
	}

	logger.Init(logger.BuildLoggers())

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
