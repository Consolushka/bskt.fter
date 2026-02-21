package main

import (
	"IMP/app/database"
	"IMP/app/internal/infra/logger"
	"IMP/app/internal/service/scheduler"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC

	if err := godotenv.Load(); err != nil {
		composite_logger.Error("Couldn't load .env file", map[string]interface{}{
			"error": err,
		})
		panic(err)
	}

	composite_logger.Init(logger.BuildSettingsFromEnv()...)

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
