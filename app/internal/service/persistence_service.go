package service

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"fmt"
	"reflect"
)

type PersistenceServiceInterface interface {
	SaveGame(game games.GameStatEntity) error
}

type PersistenceService struct {
	gamesRepo   ports.GamesRepo
	teamsRepo   ports.TeamsRepo
	playersRepo ports.PlayersRepo
}

func NewPersistenceService(gamesRepo ports.GamesRepo, teamsRepo ports.TeamsRepo, playersRepo ports.PlayersRepo) *PersistenceService {
	return &PersistenceService{
		gamesRepo:   gamesRepo,
		teamsRepo:   teamsRepo,
		playersRepo: playersRepo,
	}
}

func (s PersistenceService) SaveGame(game games.GameStatEntity) error {
	var err error

		isExists, err := s.gamesRepo.Exists(game.GameModel)

		if err != nil {

			return fmt.Errorf("Exists with %v from %s returned error: %w", game.GameModel, reflect.TypeOf(s.gamesRepo), err)

		}

	

		if isExists {

			return nil

		}

	

		game.GameModel, err = s.gamesRepo.FirstOrCreate(game.GameModel)

		if err != nil {

			return fmt.Errorf("FirstOrCreate with %v from %s returned error: %w", game.GameModel, reflect.TypeOf(s.gamesRepo), err)

		}
	game.HomeTeamStat.GameTeamStatModel.GameId = game.GameModel.Id
	game.AwayTeamStat.GameTeamStatModel.GameId = game.GameModel.Id

	err = s.saveTeamModel(&game.HomeTeamStat)
	if err != nil {
		return err
	}

	err = s.saveTeamModel(&game.AwayTeamStat)
	if err != nil {
		return err
	}

	err = s.saveTeamStatModel(&game.HomeTeamStat)
	if err != nil {
		return err
	}

	err = s.saveTeamStatModel(&game.AwayTeamStat)
	if err != nil {
		return err
	}

	for _, playerStats := range game.HomeTeamStat.PlayerStats {
		err = s.savePlayerModel(&playerStats)
		if err != nil {
			logger.Warn("savePlayerModel returned error", map[string]interface{}{
				"playerStats": playerStats,
				"error":       err,
			})
			continue
		}

		err = s.savePlayerStatModel(&playerStats)
		if err != nil {
			logger.Warn("savePlayerStatModel returned error", map[string]interface{}{
				"playerStats": playerStats,
				"error":       err,
			})
			continue
		}
	}

	for _, playerStats := range game.AwayTeamStat.PlayerStats {
		err = s.savePlayerModel(&playerStats)
		if err != nil {
			logger.Warn("savePlayerModel returned error", map[string]interface{}{
				"playerStats": playerStats,
				"error":       err,
			})
			continue
		}

		err = s.savePlayerStatModel(&playerStats)
		if err != nil {
			logger.Warn("savePlayerStatModel returned error", map[string]interface{}{
				"playerStats": playerStats,
				"error":       err,
			})
			continue
		}
	}

	return nil
}

func (s PersistenceService) saveTeamModel(teamStats *teams.TeamStatEntity) error {
	var err error

	teamStats.TeamModel, err = s.teamsRepo.FirstOrCreate(teamStats.TeamModel)
	if err != nil {
		return fmt.Errorf("FirstOrCreate with %v from %s returned error: %w", teamStats.TeamModel, reflect.TypeOf(s.teamsRepo), err)
	}
	teamStats.GameTeamStatModel.TeamId = teamStats.TeamModel.Id

	return nil
}

func (s PersistenceService) saveTeamStatModel(entity *teams.TeamStatEntity) error {
	var err error

	entity.GameTeamStatModel, err = s.teamsRepo.FirstOrCreateStats(entity.GameTeamStatModel)
	if err != nil {
		return fmt.Errorf("FirstOrCreateStats with %v from %s returned error: %w", entity.GameTeamStatModel, reflect.TypeOf(s.teamsRepo), err)
	}

	for index := range entity.PlayerStats {
		entity.PlayerStats[index].GameTeamPlayerStatModel.GameId = entity.GameTeamStatModel.GameId
		entity.PlayerStats[index].GameTeamPlayerStatModel.TeamId = entity.GameTeamStatModel.TeamId
	}

	return nil
}

func (s PersistenceService) savePlayerModel(entity *players.PlayerStatisticEntity) error {
	var err error

	entity.PlayerModel, err = s.playersRepo.FirstOrCreate(entity.PlayerModel)
	if err != nil {
		return fmt.Errorf("FirstOrCreate with %v from %s returned error: %w", entity.PlayerModel, reflect.TypeOf(s.playersRepo), err)
	}
	entity.GameTeamPlayerStatModel.PlayerId = entity.PlayerModel.Id

	return nil
}

func (s PersistenceService) savePlayerStatModel(entity *players.PlayerStatisticEntity) error {
	var err error

	entity.GameTeamPlayerStatModel, err = s.playersRepo.FirstOrCreateStat(entity.GameTeamPlayerStatModel)
	if err != nil {
		return fmt.Errorf("FirstOrCreateStat with %v from %s returned error: %w", entity.GameTeamPlayerStatModel, reflect.TypeOf(s.playersRepo), err)
	}

	return nil
}
