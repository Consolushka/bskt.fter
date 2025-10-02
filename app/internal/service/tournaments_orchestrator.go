package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"sync"
)

type TournamentsOrchestrator struct {
	persistenceService PersistenceServiceInterface
	tournamentsRepo    ports.TournamentsRepo
}

func NewTournamentsOrchestrator(persistenceService PersistenceServiceInterface, tournamentsRepo ports.TournamentsRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		persistenceService: persistenceService,
		tournamentsRepo:    tournamentsRepo,
	}

}

// ProcessAllTournamentsToday
// Fetches all active tournaments from repository and processes today games
func (t TournamentsOrchestrator) ProcessAllTournamentsToday() error {
	activeTournaments, err := t.tournamentsRepo.ListActiveTournaments()

	if err != nil {
		return err
	}

	var tournamentsGroup sync.WaitGroup
	tournamentsGroup.Add(len(activeTournaments))

	logger.Info("Start processing tournaments", map[string]interface{}{
		"tournaments": activeTournaments,
	})

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer tournamentsGroup.Done()

			statsProvider, err := NewTournamentStatsProvider(tournament)
			if err != nil {
				logger.Error("Error while creating stats provider", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
				})
				return
			}

			processor := NewTournamentProcessor(statsProvider, t.persistenceService, tournament.Id)

			logger.Info("Start processing tournament", map[string]interface{}{
				"tournament": tournament,
			})
			err = processor.Process()
			if err != nil {
				logger.Error("Error while processing tournament games", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
					"processor":  processor,
				})
				return
			}
		}(tournament)
	}

	tournamentsGroup.Wait()

	return nil
}
