package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/service"
	"IMP/app/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC

	godotenv.Load()

	logger.Init(logger.BuildLoggers())

	database.OpenDbConnection()

	db := database.GetDB()
	tournamentsOrchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
	)

	err := tournamentsOrchestrator.ProcessAllTournamentsToday()

	if err != nil {
		panic(err)
	}
}
