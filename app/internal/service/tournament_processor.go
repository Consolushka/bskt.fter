package service

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/ports"
	"fmt"
	"reflect"
	"strconv"
	"time"

	compositeLogger "github.com/Consolushka/golang.composite_logger/pkg"
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
		return fmt.Errorf("GetGamesStatsByPeriod with %v, %v from %s returned error: %w", from, to, reflect.TypeOf(t.statsProvider), err)
	}
	savedGames := make([]string, 0, len(gameEntities))

	if len(gameEntities) > 0 {
		compositeLogger.Info("Start processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
		})
	}

	for _, gameEntity := range gameEntities {
		gameEntity.GameModel.TournamentId = t.tournamentId
		isExists, err := t.gamesRepo.Exists(gameEntity.GameModel)
		if err != nil {
			compositeLogger.Error("Failed to check whether game exists", map[string]interface{}{
				"tournamentId": gameEntity.GameModel.TournamentId,
				"title":        gameEntity.GameModel.Title,
				"scheduledAt":  gameEntity.GameModel.ScheduledAt,
				"error":        err,
			})
			return fmt.Errorf("Exists with %v from %s returned error: %w", gameEntity.GameModel, reflect.TypeOf(t.gamesRepo), err)
		}
		if isExists {
			compositeLogger.Info("Game already exists. Skip game processing", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
			})
			continue
		}

		gameEntity, err = t.statsProvider.EnrichGameStats(gameEntity)
		if err != nil {
			compositeLogger.Warn("Couldn't enrich game stats", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
				"error":     err,
			})
			continue
		}

		// DISCOVERY & INGESTION for players
		t.discoverAndIngestPlayers(&gameEntity)

		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			compositeLogger.Error("t.persistenceService.SaveGame returned error", map[string]interface{}{
				"error":      err,
				"gameEntity": gameEntity,
			})
			continue
		}

		savedGames = append(savedGames, formatSavedGameLog(gameEntity))
	}

	if len(gameEntities) > 0 {
		compositeLogger.Info("Finished processing tournament games", map[string]interface{}{
			"tournamentId": t.tournamentId,
			"savedCount":   len(savedGames),
			"savedGames":   savedGames,
		})
	}

	return nil
}

func (t TournamentProcessor) discoverAndIngestPlayers(gameEntity *games.GameStatEntity) {
	homePlayers := gameEntity.HomeTeamStat.PlayerStats
	awayPlayers := gameEntity.AwayTeamStat.PlayerStats

	for i := range homePlayers {
		t.ensurePlayerBio(&homePlayers[i])
	}

	for i := range awayPlayers {
		t.ensurePlayerBio(&awayPlayers[i])
	}
}

func (t TournamentProcessor) ensurePlayerBio(playerStat *players.PlayerStatisticEntity) {
	// DISCOVERY: Check if player exists in our DB
	playersByFullName, err := t.playersRepo.ListByFullName(playerStat.PlayerModel.FullName)
	if err != nil {
		compositeLogger.Error("Failed to search players by full name", map[string]interface{}{
			"playerFullName": playerStat.PlayerModel.FullName,
			"error":          err,
		})
		return
	}

	// If player exists, we already have their bio (likely)
	if len(playersByFullName) == 1 {
		return
	}

	// INGESTION: If player is unknown or missing critical data, fetch from provider
	if playerStat.PlayerModel.FullName == "" || time.Time.IsZero(playerStat.PlayerModel.BirthDate) {
		playerBio, err := t.statsProvider.GetPlayerBio(playerStat.PlayerExternalId)
		if err != nil {
			compositeLogger.Warn("error while fetching player bio", map[string]interface{}{
				"playerId": playerStat.PlayerExternalId,
				"err":      err,
			})
		} else {
			playerStat.PlayerModel.FullName = playerBio.FullName
			playerStat.PlayerModel.BirthDate = playerBio.BirthDate
		}
	}
}

func formatSavedGameLog(gameEntity games.GameStatEntity) string {
	return gameEntity.GameModel.Title + " " +
		formatScore(gameEntity.HomeTeamStat.GameTeamStatModel.Score) + ":" +
		formatScore(gameEntity.AwayTeamStat.GameTeamStatModel.Score)
}

func formatScore(score int) string {
	return strconv.Itoa(score)
}
