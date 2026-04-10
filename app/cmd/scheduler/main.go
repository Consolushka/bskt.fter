package main

import (
	"IMP/app/database"
	"IMP/app/internal/infra/config"
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

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	composite_logger.Init(logger.BuildSettings(cfg.Logger)...)

	db := database.OpenDbConnection(cfg.Database)

	scheduler.Handle(db, cfg)
}
