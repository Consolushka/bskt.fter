package service

import (
	"IMP/app/internal/ports"
	"fmt"
	"time"
)

type TournamentProcessorInterface interface {
	Process() error
}

type TournamentProcessor struct {
	statsProvider      ports.StatsProvider
	persistenceService PersistenceServiceInterface
}

func NewTournamentProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface) *TournamentProcessor {
	return &TournamentProcessor{
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
	}
}

func (t TournamentProcessor) Process() error {
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
