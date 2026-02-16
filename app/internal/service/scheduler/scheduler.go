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
	"sync"
	"time"

	"gorm.io/gorm"
)

func Handle(db *gorm.DB) {
	pollIntervalString := os.Getenv("SCHEDULER_POLL_INTERVAL")
	pollIntervalInMinutes, err := strconv.Atoi(pollIntervalString)
	if err != nil || pollIntervalInMinutes <= 0 {
		logger.Warn("Couldn't load SCHEDULER_POLL_INTERVAL. uses default 30", map[string]interface{}{
			"value": pollIntervalString,
			"error": err,
		})
		pollIntervalInMinutes = 30
	}

	logger.Info("Scheduler started", map[string]interface{}{
		"pollIntervalInMinutes": pollIntervalInMinutes,
	})

	executePollingCycle(db)

	ticker := time.NewTicker(time.Duration(pollIntervalInMinutes) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		executePollingCycle(db)
	}
}

func executePollingCycle(db *gorm.DB) {
	var watermarkRepo ports.PollWatermarkRepo = poll_watermarks_repo.NewGormRepo(db)
	var tournamentsRepo ports.TournamentsRepo = tournaments_repo.NewGormRepo(db)

	orchestrator := service.NewTournamentsOrchestrator(
		*service.NewPersistenceService(games_repo.NewGormRepo(db), teams_repo.NewGormRepo(db), players_repo.NewGormRepo(db)),
		tournamentsRepo,
		players_repo.NewGormRepo(db),
		games_repo.NewGormRepo(db),
	)

	activeTournaments, err := tournamentsRepo.ListActive()
	if err != nil {
		logger.Error("Couldn't fetch active tournaments", map[string]interface{}{
			"error": err,
		})
		return
	}

	now := time.Now().UTC()
	var wg sync.WaitGroup
	wg.Add(len(activeTournaments))

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer wg.Done()

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
		}(tournament)
	}

	wg.Wait()

	logger.Info("Polling cycle finished", map[string]interface{}{
		"finishedAt": now,
	})
}

func toStartOfUTCDay(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}

