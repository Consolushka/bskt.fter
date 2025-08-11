package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"fmt"
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
// ProcessAllTournamentsToday
// Fetches all active tournaments from repository and processes today games
func (t TournamentsOrchestrator) ProcessAllTournamentsToday() error {
	activeTournaments, err := t.tournamentsRepo.ListActiveTournaments()

	if err != nil {
		return err
	}

	var tournamentsGroup sync.WaitGroup
	tournamentsGroup.Add(len(activeTournaments))

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer tournamentsGroup.Done()

			statsProvider, err := NewStatsProvider(tournament.League.Alias)
			if err != nil {
				fmt.Println("There was an error creating stats provider. Error: ", err)
				return
			}

			processor := NewTournamentsProcessor(statsProvider, t.persistenceService)
			err = processor.Process(tournament)
			if err != nil {
				fmt.Println("There was an error processing tournament games. Error: ", err)
				return
			}
		}(tournament)
	}

	tournamentsGroup.Wait()

	fmt.Println("everything is done")

	return nil
}
