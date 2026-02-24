package api_nba

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type EntityTransformer struct {
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

func (e *EntityTransformer) MapPlayerStatistics(response PlayerStatisticResponse, homeTeamId int, awayTeamId int, gameBusinessEntity *games.GameStatEntity) error {
	homeTeamPlayers := make([]players.PlayerStatisticEntity, 0)
	awayTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	for _, playerStat := range response.Response {
		playerStatEntity := players.PlayerStatisticEntity{}

		playerStatsErr := e.mapPlayerStatistic(playerStat, &playerStatEntity)
		if playerStatsErr != nil {
			return fmt.Errorf("error mapping player statistic for player %d: %w", playerStat.Player.Id, playerStatsErr)
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

func (e *EntityTransformer) mapPlayerStatistic(player PlayerStatisticEntity, playerBusinessEntity *players.PlayerStatisticEntity) error {
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

	fgp, err := strconv.ParseFloat(player.Fgp, 32)
	if err != nil {
		return errors.New(err.Error())
	}

	*playerBusinessEntity = players.PlayerStatisticEntity{
		PlayerExternalId: strconv.Itoa(player.Player.Id),
		PlayerModel: players.PlayerModel{
			FullName: player.Player.Firstname + " " + player.Player.Lastname,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds:        secondsPlayed,
			PlsMin:               int8(plsMin),
			Points:               uint8(player.Points),
			Rebounds:             uint8(player.TotReb),
			Assists:              uint8(player.Assists),
			Steals:               uint8(player.Steals),
			Blocks:               uint8(player.Blocks),
			FieldGoalsPercentage: float32(fgp) / 100,
			Turnovers:            uint8(player.Turnovers),
		},
	}

	return nil
}

func NewEntityTransformer() EntityTransformer {
	return EntityTransformer{}
}
