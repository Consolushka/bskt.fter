package scheduler

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/poll_watermarks_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/poll_watermarks"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/pkg/logger"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func Handle(db *gorm.DB) {
	pollIntervalInMinutes := getEnvInt("SCHEDULER_POLL_INTERVAL", 30)
	staggerIntervalInMinutes := getEnvInt("SCHEDULER_STAGGER_INTERVAL_MINUTES", 5)

	logger.Info("Scheduler starting staggered workers", map[string]interface{}{
		"pollIntervalInMinutes":    pollIntervalInMinutes,
		"staggerIntervalInMinutes": staggerIntervalInMinutes,
	})

	tournamentsRepo := tournaments_repo.NewGormRepo(db)
	activeTournaments, err := tournamentsRepo.ListActive()
	if err != nil {
		logger.Error("Couldn't fetch active tournaments", map[string]interface{}{
			"error": err,
		})
		return
	}

	for _, tournament := range activeTournaments {
		go runTournamentWorker(db, tournament, time.Duration(pollIntervalInMinutes)*time.Minute)

		// Wait before starting the next worker to distribute the load
		time.Sleep(time.Duration(staggerIntervalInMinutes) * time.Minute)
	}

	// Keep the main goroutine alive
	select {}
}

func runTournamentWorker(db *gorm.DB, tournament tournaments.TournamentModel, interval time.Duration) {
	logger.Info("Worker started", map[string]interface{}{
		"tournamentId": tournament.Id,
		"interval":     interval.String(),
	})

	// Immediate first run
	processTournament(db, tournament)

	// Periodic runs
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		processTournament(db, tournament)
	}
}

func processTournament(db *gorm.DB, tournament tournaments.TournamentModel) {
	var watermarkRepo ports.PollWatermarkRepo = poll_watermarks_repo.NewGormRepo(db)
	var tournamentsRepo ports.TournamentsRepo = tournaments_repo.NewGormRepo(db)

	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournamentsRepo,
		players_repo.NewGormRepo(db),
		games_repo.NewGormRepo(db),
	)

	now := time.Now().UTC()
	startOfDay := toStartOfUTCDay(now)

	watermarkModel, err := watermarkRepo.FirstOrCreate(poll_watermarks.PollWatermarkModel{
		TournamentId:         tournament.Id,
		LastSuccessfulPollAt: startOfDay,
	})
	if err != nil {
		logger.Error("Couldn't read or create tournament watermark", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        err,
		})
		return
	}

	oldPollAt := watermarkModel.LastSuccessfulPollAt
	if err = orchestrator.ProcessTournament(tournament, oldPollAt, now); err != nil {
		logger.Error("Error while processing tournament games", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        err,
		})
		return
	}

	watermarkModel.LastSuccessfulPollAt = now
	_, err = watermarkRepo.Update(watermarkModel)
	if err != nil {
		logger.Warn("Couldn't update tournament watermark", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        err,
		})
		return
	}

	logger.Info("Tournament processed successfully", map[string]interface{}{
		"tournamentId": tournament.Id,
		"from":         oldPollAt,
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
