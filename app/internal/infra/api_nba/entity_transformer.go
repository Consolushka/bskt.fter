package api_nba

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/logger"
	"errors"
	"strconv"
	"strings"
	"time"
)

type EntityTransformer struct {
	client ClientInterface
}

func (e *EntityTransformer) Transform(game GameEntity) (games.GameStatEntity, error) {
	businessEntity := games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: game.Date.Start,
			Title:       game.Teams.Home.Code + " - " + game.Teams.Visitors.Code,
		},
		HomeTeamStat: teams.TeamStatEntity{
			TeamModel: teams.TeamModel{
				Name:     game.Teams.Home.Nickname,
				HomeTown: strings.TrimRight(strings.Replace(game.Teams.Home.Name, game.Teams.Home.Nickname, "", 1), " "),
			},
			GameTeamStatModel: teams.GameTeamStatModel{
				Score:     game.Scores.Home.Points,
				FinalDiff: game.Scores.Home.Points - game.Scores.Visitors.Points,
			},
			PlayerStats: nil,
		},
		AwayTeamStat: teams.TeamStatEntity{
			TeamModel: teams.TeamModel{
				Name:     game.Teams.Visitors.Nickname,
				HomeTown: strings.TrimRight(strings.Replace(game.Teams.Visitors.Name, game.Teams.Visitors.Nickname, "", 1), " "),
			},
			GameTeamStatModel: teams.GameTeamStatModel{
				Score:     game.Scores.Visitors.Points,
				FinalDiff: game.Scores.Visitors.Points - game.Scores.Home.Points,
			},
			PlayerStats: nil,
		},
	}

	err := e.enrichGamePlayers(game, &businessEntity)
	if err != nil {
		return games.GameStatEntity{}, err
	}

	return businessEntity, nil
}

func (e *EntityTransformer) enrichGamePlayers(game GameEntity, gameBusinessEntity *games.GameStatEntity) error {
	homeTeamId := game.Teams.Home.Id
	homeTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	awayTeamId := game.Teams.Visitors.Id
	awayTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	gameStat, err := e.client.PlayersStatistics(0, game.Id, 0, "")
	if err != nil {
		return err
	}

	for _, playerStat := range gameStat.Response {
		playerStatEntity := players.PlayerStatisticEntity{}

		playerStatsErr := e.enrichPlayerStatistic(playerStat, &playerStatEntity)
		if playerStatsErr != nil {
			logger.Warn("There was an error with player statistics", map[string]interface{}{
				"playerStat":       playerStat,
				"playerStatEntity": playerStatEntity,
				"error":            playerStatsErr,
			})
		}

		if playerStat.Team.Id == homeTeamId {
			homeTeamPlayers = append(homeTeamPlayers, playerStatEntity)
		} else if playerStat.Team.Id == awayTeamId {
			awayTeamPlayers = append(awayTeamPlayers, playerStatEntity)
		}
	}

	gameBusinessEntity.HomeTeamStat.PlayerStats = homeTeamPlayers
	gameBusinessEntity.AwayTeamStat.PlayerStats = awayTeamPlayers

	return nil
}

func (e *EntityTransformer) enrichPlayerStatistic(player PlayerStatisticEntity, playerBusinessEntity *players.PlayerStatisticEntity) error {
	splittedGameTime := strings.Split(player.Min, ":")
	minutesPlayed, err := strconv.Atoi(splittedGameTime[0])
	if err != nil {
		return errors.New(err.Error())
	}

	var secondsAfterMinutes int
	if len(splittedGameTime) == 1 {
		secondsAfterMinutes = 0
	} else {
		secondsAfterMinutes, err = strconv.Atoi(splittedGameTime[1])
		if err != nil {
			return errors.New(err.Error())
		}
	}

	secondsPlayed := minutesPlayed*60 + secondsAfterMinutes
	plsMin, err := strconv.Atoi(player.PlusMinus)
	if err != nil {
		return errors.New(err.Error())
	}

	*playerBusinessEntity = players.PlayerStatisticEntity{
		PlayerModel: players.PlayerModel{
			FullName: player.Player.Firstname + " " + player.Player.Lastname,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds: secondsPlayed,
			PlsMin:        int8(plsMin),
		},
	}

	playerBioErr := e.enrichPlayerBio(player.Player.Id, playerBusinessEntity)
	if playerBioErr != nil {
		return playerBioErr
	}

	return nil
}

func (e *EntityTransformer) enrichPlayerBio(playerId int, playerEntity *players.PlayerStatisticEntity) error {
	playerBio, err := e.client.PlayerInfo(playerId, "", 0, 0, "", "")
	if err != nil {
		return err
	}

	// For free plan limit is 10 requests/minute
	//time.Sleep(6 * time.Second)
	// todo: some players doesn't have birthdate
	playerEntity.PlayerModel.BirthDate, err = time.Parse("2006-01-02", playerBio.Response[0].Birth.Date)
	if err != nil {
		playerEntity.PlayerModel.BirthDate = time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
		return errors.New(err.Error())
	}

	return nil
}

func NewEntityTransformer(client ClientInterface) EntityTransformer {
	return EntityTransformer{
		client: client,
	}
}
