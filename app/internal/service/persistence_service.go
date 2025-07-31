package service

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/internal/ports"
	"fmt"
)

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

	game.GameModel, err = s.gamesRepo.FindOrCreateGame(game.GameModel)
	if err != nil {
		return err
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
			// todo: log
			fmt.Println("There was an error. Error: ", err)
			continue
		}

		err = s.savePlayerStatModel(&playerStats)
		if err != nil {
			// todo: log
			fmt.Println("There was an error. Error: ", err)
			continue
		}
	}

	for _, playerStats := range game.AwayTeamStat.PlayerStats {
		err = s.savePlayerModel(&playerStats)
		if err != nil {
			// todo: log
			fmt.Println("There was an error. Error: ", err)
			continue
		}

		err = s.savePlayerStatModel(&playerStats)
		if err != nil {
			// todo: log
			fmt.Println("There was an error. Error: ", err)
			continue
		}
	}

	return nil
}

func (s PersistenceService) saveTeamModel(teamStats *teams.TeamStatEntity) error {
	var err error

	teamStats.TeamModel, err = s.teamsRepo.FirstOrCreateTeam(teamStats.TeamModel)
	if err != nil {
		return err
	}
	teamStats.GameTeamStatModel.TeamId = teamStats.TeamModel.Id

	return nil
}

func (s PersistenceService) saveTeamStatModel(entity *teams.TeamStatEntity) error {
	var err error

	entity.GameTeamStatModel, err = s.teamsRepo.FirstOrCreateTeamStats(entity.GameTeamStatModel)
	if err != nil {
		return err
	}

	for _, playerStat := range entity.PlayerStats {
		playerStat.GameTeamPlayerStatModel.GameTeamId = entity.GameTeamStatModel.Id
	}

	return nil
}

func (s PersistenceService) savePlayerModel(entity *players.PlayerStatisticEntity) error {
	var err error

	entity.PlayerModel, err = s.playersRepo.FirstOrCreatePlayer(entity.PlayerModel)
	if err != nil {
		return err
	}
	entity.GameTeamPlayerStatModel.PlayerId = entity.PlayerModel.Id

	return nil
}

func (s PersistenceService) savePlayerStatModel(entity *players.PlayerStatisticEntity) error {
	var err error

	entity.GameTeamPlayerStatModel, err = s.playersRepo.FirstOrCreatePlayerStat(entity.GameTeamPlayerStatModel)
	if err != nil {
		return err
	}

	return nil
}
