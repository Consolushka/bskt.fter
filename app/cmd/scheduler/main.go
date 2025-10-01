package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/service"
	"IMP/app/pkg/logger"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	loc, _ := time.LoadLocation("Europe/Moscow")
	time.Local = loc

	godotenv.Load()

	logger.Init(logger.BuildLoggers())

	database.OpenDbConnection()

	scheduleTasks()
}

func scheduleTasks() {
	db := database.GetDB()

	tournamentsOrchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournaments_repo.NewGormRepo(db),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go scheduleDailyTask(tournamentsOrchestrator.ProcessAllTournamentsToday, &wg, ctx, 10, 38, 0)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
	wg.Wait()
	log.Println("Stopped.")
}

func scheduleDailyTask(task func() error, wg *sync.WaitGroup, context context.Context, hour int, minute int, second int) {
	defer wg.Done()

	for {
		now := time.Now()

		target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, now.Nanosecond(), now.Location())

		var sleepDuration time.Duration

		if now.Before(target) {
			sleepDuration = target.Sub(now)
		} else {
			sleepDuration = target.Add(time.Hour * 24).Sub(now)
		}

		fmt.Println(sleepDuration.String())

		timer := time.NewTimer(sleepDuration)

		select {
		case <-context.Done():
			return
		case <-timer.C:
			go func() {
				err := task()

				if err != nil {
					logger.Error("Error while processing tournament games", map[string]interface{}{
						"error": err,
					})
				}
			}()
		}

	}
}
