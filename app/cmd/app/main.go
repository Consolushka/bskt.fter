package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/service"
	"github.com/joho/godotenv"
	"time"
)

func main() {
	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc

	err := godotenv.Load()
	if err != nil {
		panic("Couldn't load env file")
	}

	database.OpenDbConnection()

	db := database.GetDB()
	tournamentsOrchestrator := service.NewTournamentsOrchestrator(
		games_repo.NewGormRepo(db),
		teams_repo.NewGormRepo(db),
		players_repo.NewGormRepo(db),
		tournaments_repo.NewGormRepo(db),
	)

	err = tournamentsOrchestrator.ProcessAllTournamentsToday()

	if err != nil {
		panic(err)
	}
}
