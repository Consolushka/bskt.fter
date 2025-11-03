package service

import (
	"IMP/app/internal/core/players"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"time"
)

type TournamentProcessorInterface interface {
	ProcessByPeriod(from, to time.Time) error
}

type TournamentProcessor struct {
	tournamentId       uint
	statsProvider      ports.StatsProvider
	persistenceService PersistenceServiceInterface
	playersRepo        ports.PlayersRepo
}

func NewTournamentProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface, playersRepo ports.PlayersRepo, tournamentId uint) *TournamentProcessor {
	return &TournamentProcessor{
		tournamentId:       tournamentId,
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
		playersRepo:        playersRepo,
	}
}

func (t TournamentProcessor) ProcessByPeriod(from, to time.Time) error {
	gameEntities, err := t.statsProvider.GetGamesStatsByPeriod(from, to)
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

		var allPlayers []players.PlayerStatisticEntity

		allPlayers = append(allPlayers, gameEntity.HomeTeamStat.PlayerStats...)
		allPlayers = append(allPlayers, gameEntity.AwayTeamStat.PlayerStats...)

		for _, playerStat := range allPlayers {
			playersByFullName, err := t.playersRepo.PlayersByFullName(playerStat.PlayerModel.FullName)
			if err != nil {
				return err
			}

			if len(playersByFullName) != 1 {
				logger.Info("Player not found in database", map[string]interface{}{
					"playerFullName": playerStat.PlayerModel.FullName,
				})

				if playerStat.PlayerModel.FullName == "" || time.Time.IsZero(playerStat.PlayerModel.BirthDate) {
					playerBio, err := t.statsProvider.GetPlayerBio(playerStat.PlayerExternalId)
					if err != nil {
						logger.Warn("error while fetching player bio", map[string]interface{}{
							"err": err,
						})
					} else {
						playerStat.PlayerModel.FullName = playerBio.FullName
						playerStat.PlayerModel.BirthDate = playerBio.BirthDate
					}
				}
			}
		}

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
