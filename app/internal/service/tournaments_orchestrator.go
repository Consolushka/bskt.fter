package service

import (
	"IMP/app/internal/ports"
	"fmt"
	"sync"
	"time"
)

type TournamentsOrchestrator struct {
	gamesRepo       ports.GamesRepo
	teamsRepo       ports.TeamsRepo
	playersRepo     ports.PlayersRepo
	tournamentsRepo ports.TournamentsRepo
}

func NewTournamentsOrchestrator(gamesRepo ports.GamesRepo, teamsRepo ports.TeamsRepo, playersRepo ports.PlayersRepo, tournamentsRepo ports.TournamentsRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		gamesRepo:       gamesRepo,
		teamsRepo:       teamsRepo,
		playersRepo:     playersRepo,
		tournamentsRepo: tournamentsRepo,
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

	persistenceService := NewPersistenceService(t.gamesRepo, t.teamsRepo, t.playersRepo)

	for _, tournament := range activeTournaments {
		//todo: когда добавлю запись в бд, то сюда добавить каналы
		go func() {
			statsProvider, err := NewStatsProvider(tournament.League.Alias)
			if err != nil {
				fmt.Println("There was an error. Error: ", err)
				tournamentsGroup.Done()
				return
			}

			gameEntities, err := statsProvider.GetGamesStatsByDate(time.Now())
			if err != nil {
				fmt.Println("There was an error. Error: ", err)
			}

			for _, gameEntity := range gameEntities {
				err = persistenceService.SaveGame(gameEntity)
				if err != nil {
					fmt.Println("Error while saving game to db: ", err)
					continue
				}
			}

			tournamentsGroup.Done()
		}()
	}

	tournamentsGroup.Wait()

	fmt.Println("everything is done")

	return nil
}
