package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/executable_by_scheduler"
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/pkg/logger"
	"IMP/app/pkg/time_utils"
	"context"
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
	go scheduleDailyTask(executable_by_scheduler.NewProcessTodayTournaments(tournamentsOrchestrator), &wg, ctx)

	wg.Add(1)
	go scheduleDailyTask(executable_by_scheduler.NewProcessNightlyTournaments(tournamentsOrchestrator), &wg, ctx)

	wg.Add(1)
	go scheduleDailyTask(executable_by_scheduler.NewProcessYesterdayDailyTournaments(tournamentsOrchestrator), &wg, ctx)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	cancel()
	wg.Wait()
	log.Println("Stopped.")
}

func scheduleDailyTask(task ports.ExecutableByScheduler, wg *sync.WaitGroup, context context.Context) {
	defer wg.Done()

	for {
		now := time.Now()

		target := time_utils.ToMoscowTZ(task.ShouldBeExecutedAt())

		var sleepDuration time.Duration

		if now.Before(target) {
			sleepDuration = target.Sub(now)

			logger.Info(task.GetName()+" will be executed at "+target.Format("02-01-2006 15:04"), map[string]interface{}{})
		} else {
			logger.Info(task.GetName()+" should been executed at "+target.Format("02-01-2006 15:04")+". Executing...", map[string]interface{}{})

			err := task.Execute()

			if err != nil {
				logger.Error("Error while processing tournament games", map[string]interface{}{
					"error": err,
				})
			}

			target = target.Add(time.Hour * 24)

			logger.Info(task.GetName()+" will be executed at "+target.Format("02-01-2006 15:04"), map[string]interface{}{})
			sleepDuration = target.Sub(now)
		}

		timer := time.NewTimer(sleepDuration)

		select {
		case <-context.Done():
			return
		case <-timer.C:
			go func() {
				err := task.Execute()

				if err != nil {
					logger.Error("Error while processing tournament games", map[string]interface{}{
						"error": err,
					})
				}

				wg.Done()
			}()
		}

	}
}
