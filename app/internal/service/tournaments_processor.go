package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"fmt"
	"time"
)

type TournamentsProcessorInterface interface {
	Process(model tournaments.TournamentModel) error
}

type TournamentsProcessor struct {
	statsProvider      ports.StatsProvider
	persistenceService PersistenceServiceInterface
}

func NewTournamentsProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface) *TournamentsProcessor {
	return &TournamentsProcessor{
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
	}
}

func (t TournamentsProcessor) Process(model tournaments.TournamentModel) error {
	gameEntities, err := t.statsProvider.GetGamesStatsByDate(time.Now())
	if err != nil {
		return err
	}

	for _, gameEntity := range gameEntities {
		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			fmt.Println("Error while saving game to db: ", err)
			continue
		}
	}

	return nil
}
