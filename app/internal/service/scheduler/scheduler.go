package scheduler

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournament_poll_logs_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"os"
	"strconv"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"gorm.io/gorm"
)

func Handle(db *gorm.DB) {
	pollIntervalInMinutes := getEnvInt("SCHEDULER_POLL_INTERVAL", 30)
	staggerIntervalInMinutes := getEnvInt("SCHEDULER_STAGGER_INTERVAL_MINUTES", 5)

	composite_logger.Info("Scheduler starting staggered workers", map[string]interface{}{
		"pollIntervalInMinutes":    pollIntervalInMinutes,
		"staggerIntervalInMinutes": staggerIntervalInMinutes,
	})

	tournamentsRepo := tournaments_repo.NewGormRepo(db)
	gamesRepo := games_repo.NewGormRepo(db)
	teamsRepo := teams_repo.NewGormRepo(db)
	playersRepo := players_repo.NewGormRepo(db)
	pollLogRepo := tournament_poll_logs_repo.NewGormRepo(db)

	persistenceService := service.NewPersistenceService(gamesRepo, teamsRepo, playersRepo)
	orchestrator := service.NewTournamentsOrchestrator(
		persistenceService,
		tournamentsRepo,
		playersRepo,
		gamesRepo,
		pollLogRepo,
	)

	activeTournaments, err := tournamentsRepo.ListActive()
	if err != nil {
		composite_logger.Error("Couldn't fetch active tournaments", map[string]interface{}{
			"error": err,
		})
		return
	}

	for _, tournament := range activeTournaments {
		go runTournamentWorker(orchestrator, pollLogRepo, tournament, time.Duration(pollIntervalInMinutes)*time.Minute)

		// Wait before starting the next worker to distribute the load
		time.Sleep(time.Duration(staggerIntervalInMinutes) * time.Minute)
	}

	// Keep the main goroutine alive
	select {}
}

func runTournamentWorker(
	orchestrator *service.TournamentsOrchestrator,
	pollLogRepo ports.TournamentPollLogsRepo,
	tournament tournaments.TournamentModel,
	interval time.Duration,
) {
	defer composite_logger.Recover(map[string]interface{}{
		"tournamentId": tournament.Id,
	})

	composite_logger.Info("Worker started", map[string]interface{}{
		"tournamentId": tournament.Id,
		"interval":     interval.String(),
	})

	// Immediate first run
	processTournament(orchestrator, pollLogRepo, tournament)

	// Periodic runs
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		processTournament(orchestrator, pollLogRepo, tournament)
	}
}

func processTournament(
	orchestrator *service.TournamentsOrchestrator,
	pollLogRepo ports.TournamentPollLogsRepo,
	tournament tournaments.TournamentModel,
) {
	now := time.Now().UTC()

	// 1. DISCOVERY: Get latest successful poll interval end
	latestLog, err := pollLogRepo.GetLatestSuccess(tournament.Id)
	var intervalStart time.Time
	if err != nil {
		// If no logs found, start from today
		intervalStart = toStartOfUTCDay(now)
	} else {
		intervalStart = latestLog.IntervalEnd
	}

	// 2. INGESTION: Run orchestration (it will handle internal poll logging)
	if err := orchestrator.ProcessTournament(tournament, intervalStart, now); err != nil {
		composite_logger.Error("Error while processing tournament games", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        err,
		})
		return
	}

	composite_logger.Info("Tournament worker cycle finished", map[string]interface{}{
		"tournamentId": tournament.Id,
		"from":         intervalStart,
		"to":           now,
	})
}

func toStartOfUTCDay(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

func getEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
