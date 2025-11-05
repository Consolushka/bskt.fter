package main

import (
	"IMP/app/database"
	"IMP/app/pkg/logger"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	time.Local = time.UTC

	godotenv.Load()

	logger.Init(logger.BuildLoggers())

	database.OpenDbConnection()

	database.GetDB()
	//tournamentsOrchestrator := service.NewTournamentsOrchestrator(
	//	*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
	//	tournaments_repo.NewGormRepo(db),
	//	players_repo.NewGormRepo(db),
	//)
	//
	//err := tournamentsOrchestrator.ProcessNotUrgentEuropeanTournaments(time.Date(2025, 11, 04, 4, 59, 59, 59, time.UTC), time.Date(2025, 11, 05, 4, 59, 59, 59, time.UTC))
	//
	//if err != nil {
	//	panic(err)
	//}
}
