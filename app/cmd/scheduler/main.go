package main

import (
	"IMP/app/database"
	"IMP/app/internal/service/scheduler"
	"IMP/app/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc

	godotenv.Load()

	logger.Init(logger.BuildLoggers())

	db := database.OpenDbConnection()

	scheduler.Handle(db)
}
