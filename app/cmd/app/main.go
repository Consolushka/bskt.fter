package main

import (
	"IMP/app/database"
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

	tournamentsOrchestrator := service.NewTournamentsOrchestrator(tournaments_repo.NewGormRepo(database.GetDB()))

	err = tournamentsOrchestrator.ProcessAllTournamentsToday()

	if err != nil {
		panic(err)
	}
}
