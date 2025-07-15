package service

import (
	"IMP/app/internal/ports"
	"fmt"
	"sync"
	"time"
)

type TournamentsOrchestrator struct {
	repo ports.TournamentsRepo
}

func NewTournamentsOrchestrator(repo ports.TournamentsRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		repo: repo,
	}

}

// ProcessAllTournamentsToday
// Fetches all active tournaments from repository and processes today games
func (t TournamentsOrchestrator) ProcessAllTournamentsToday() error {
	activeTournaments, err := t.repo.ListActiveTournaments()

	if err != nil {
		return err
	}

	var tournamentsGroup sync.WaitGroup
	tournamentsGroup.Add(len(activeTournaments))

	for _, tournament := range activeTournaments {
		go func() {
			statsProvider, err := NewStatsProvider(tournament.League.Alias)
			if err != nil {
				fmt.Println("There was an error. Error: ", err)
			}

			_, err = statsProvider.GetGamesStatsByDate(time.Now())
			if err != nil {
				fmt.Println("There was an error. Error: ", err)
			}

			tournamentsGroup.Done()
		}()
	}

	tournamentsGroup.Wait()

	fmt.Println("everything is done")

	return nil
}
