package service

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"strconv"
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
	gamesRepo          ports.GamesRepo
}

func NewTournamentProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface, playersRepo ports.PlayersRepo, gamesRepo ports.GamesRepo, tournamentId uint) *TournamentProcessor {
	return &TournamentProcessor{
		tournamentId:       tournamentId,
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
		playersRepo:        playersRepo,
		gamesRepo:          gamesRepo,
	}
}

func (t TournamentProcessor) ProcessByPeriod(from, to time.Time) error {
	gameEntities, err := t.statsProvider.GetGamesStatsByPeriod(from, to)
	if err != nil {
		return err
	}
	savedGames := make([]string, 0, len(gameEntities))

	if len(gameEntities) > 0 {
		logger.Info("Start processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
		})
	}

	for _, gameEntity := range gameEntities {
		gameEntity.GameModel.TournamentId = t.tournamentId
		isExists, err := t.gamesRepo.GameExists(gameEntity.GameModel)
		if err != nil {
			return err
		}
		if isExists {
			logger.Info("Game already exists. Skip game processing", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
			})
			continue
		}

		gameEntity, err = t.statsProvider.EnrichGameStats(gameEntity)
		if err != nil {
			logger.Warn("Couldn't enrich game stats", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
				"error":     err,
			})
			continue
		}

		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			logger.Error("t.persistenceService.SaveGame returned error", map[string]interface{}{
				"error":      err,
				"gameEntity": gameEntity,
			})
			continue
		}

		var allPlayers []players.PlayerStatisticEntity

		allPlayers = append(allPlayers, gameEntity.HomeTeamStat.PlayerStats...)
		allPlayers = append(allPlayers, gameEntity.AwayTeamStat.PlayerStats...)

		for _, playerStat := range allPlayers {
			playersByFullName, err := t.playersRepo.PlayersByFullName(playerStat.PlayerModel.FullName)
			if err != nil {
				return err
			}

			if len(playersByFullName) != 1 {
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
		savedGames = append(savedGames, formatSavedGameLog(gameEntity))
	}

	if len(gameEntities) > 0 {
		logger.Info("Finished processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
			"savedCount":   len(savedGames),
			"savedGames":   savedGames,
		})
	}

	return nil
}

func formatSavedGameLog(gameEntity games.GameStatEntity) string {
	return gameEntity.GameModel.Title + " " +
		formatScore(gameEntity.HomeTeamStat.GameTeamStatModel.Score) + ":" +
		formatScore(gameEntity.AwayTeamStat.GameTeamStatModel.Score)
}

func formatScore(score int) string {
	return strconv.Itoa(score)
}
