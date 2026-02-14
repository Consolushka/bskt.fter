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
	businessEntity := e.TransformWithoutPlayers(game)

	err := e.EnrichGamePlayers(game.Id, game.Teams.Home.Id, game.Teams.Visitors.Id, &businessEntity)
	if err != nil {
		return games.GameStatEntity{}, err
	}

	return businessEntity, nil
}

func (e *EntityTransformer) TransformWithoutPlayers(game GameEntity) games.GameStatEntity {
	duration := 0
	for i := 0; i < game.Periods.Total; i++ {
		if i <= 4 {
			duration += 12
		} else {
			duration += 5
		}
	}

	businessEntity := games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: game.Date.Start,
			Duration:    duration,
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

	return businessEntity
}

func (e *EntityTransformer) EnrichGamePlayers(gameId int, homeTeamId int, awayTeamId int, gameBusinessEntity *games.GameStatEntity) error {
	homeTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	awayTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	gameStat, err := e.client.PlayersStatistics(0, gameId, 0, "")
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
			continue
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

func (e *EntityTransformer) enrichGamePlayers(game GameEntity, gameBusinessEntity *games.GameStatEntity) error {
	return e.EnrichGamePlayers(game.Id, game.Teams.Home.Id, game.Teams.Visitors.Id, gameBusinessEntity)
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
		PlayerExternalId: strconv.Itoa(player.Player.Id),
		PlayerModel: players.PlayerModel{
			FullName: player.Player.Firstname + " " + player.Player.Lastname,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds: secondsPlayed,
			PlsMin:        int8(plsMin),
		},
	}

	err = e.enrichPlayerBio(player.Player.Id, playerBusinessEntity)
	if err != nil {
		return err
	}

	return nil
}

func (e *EntityTransformer) enrichPlayerBio(playerId int, playerBusinessEntity *players.PlayerStatisticEntity) error {
	playerResponse, err := e.client.PlayerInfo(playerId, "", 0, 0, "", "")
	if err != nil {
		return err
	}

	if len(playerResponse.Response) == 0 {
		return errors.New("empty player info response")
	}

	birthDate, err := time.Parse("2006-01-02", playerResponse.Response[0].Birth.Date)
	if err != nil {
		return err
	}

	playerBusinessEntity.PlayerModel.BirthDate = birthDate
	return nil
}

func NewEntityTransformer(client ClientInterface) EntityTransformer {
	return EntityTransformer{
		client: client,
	}
}
