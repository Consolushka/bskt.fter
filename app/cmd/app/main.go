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
	//)
	//
	//err := tournamentsOrchestrator.ProcessUrgentEuropeanTournaments(time.Date(2025, 10, 9, 16, 00, 00, 00, time.UTC), time.Now())
	//
	//if err != nil {
	//	panic(err)
	//}
}
