package service

import (
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"time"
)

type TournamentProcessorInterface interface {
	Process(date time.Time) error
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

func (t TournamentProcessor) Process(date time.Time) error {
	gameEntities, err := t.statsProvider.GetGamesStatsByDate(date)
	if err != nil {
		return err
	}

	if len(gameEntities) > 0 {
		logger.Info("Start processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
		})
	}

	for _, gameEntity := range gameEntities {
		gameEntity.GameModel.TournamentId = t.tournamentId
		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			logger.Error("t.persistenceService.SaveGame returned error", map[string]interface{}{
				"error":      err,
				"gameEntity": gameEntity,
			})
			continue
		}

		logger.Info("Game was successfully saved", map[string]interface{}{
			"gameEntity": gameEntity,
		})
	}

	if len(gameEntities) > 0 {
		logger.Info("Finished processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
		})
	}

	return nil
}
