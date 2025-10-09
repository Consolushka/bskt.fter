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

	godotenv.Load()

	logger.Init(logger.BuildLoggers())

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
