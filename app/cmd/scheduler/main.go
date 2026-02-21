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

	//nolint:errcheck
	godotenv.Load()

	composite_logger.Init(logger.BuildSettingsFromEnv()...)

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
