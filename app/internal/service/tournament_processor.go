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
	tournamentId       uint
	statsProvider      ports.StatsProvider
	persistenceService PersistenceServiceInterface
}

func NewTournamentProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface, tournamentId uint) *TournamentProcessor {
	return &TournamentProcessor{
		tournamentId:       tournamentId,
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
	}
}

func (t TournamentProcessor) Process() error {
	gameEntities, err := t.statsProvider.GetGamesStatsByDate(time.Now().Add(-time.Hour * 24))
	if err != nil {
		return err
	}

	for _, gameEntity := range gameEntities {
		gameEntity.GameModel.TournamentId = t.tournamentId
		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			fmt.Println("Error while saving game to db: ", err)
			continue
		}
	}

	return nil
}
