package api_basketball

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/statsutil"
	"fmt"
	"strconv"
	"strings"
)

type EntityTransformer struct {
}

func (e *EntityTransformer) TransformWithoutPlayers(game GameEntity) games.GameStatEntity {
	// Assume 12 mins for quarter by default if it's NBA, or maybe just look at the quarters
	// api-basketball doesn't explicitly tell quarter duration, but we can assume standard 4 quarters + overtimes
	duration := 48 // 4 * 12
	if game.Scores.Home.OverTime != nil || game.Scores.Away.OverTime != nil {
		duration += 5 // standard overtime
	}

	businessEntity := games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: game.Date,
			Duration:    duration,
			Title:       game.Teams.Home.Name + " - " + game.Teams.Away.Name,
		},
		HomeTeamStat: teams.TeamStatEntity{
			TeamModel: teams.TeamModel{
				Name: game.Teams.Home.Name,
			},
			GameTeamStatModel: teams.GameTeamStatModel{
				Score:     game.Scores.Home.Total,
				FinalDiff: game.Scores.Home.Total - game.Scores.Away.Total,
			},
			PlayerStats: nil,
		},
		AwayTeamStat: teams.TeamStatEntity{
			TeamModel: teams.TeamModel{
				Name: game.Teams.Away.Name,
			},
			GameTeamStatModel: teams.GameTeamStatModel{
				Score:     game.Scores.Away.Total,
				FinalDiff: game.Scores.Away.Total - game.Scores.Home.Total,
			},
			PlayerStats: nil,
		},
	}

	return businessEntity
}

func (e *EntityTransformer) MapPlayerStatistics(response PlayerStatsResponse, homeTeamId int, awayTeamId int, gameBusinessEntity *games.GameStatEntity) error {
	homeTeamPlayers := make([]players.PlayerStatisticEntity, 0)
	awayTeamPlayers := make([]players.PlayerStatisticEntity, 0)

	for _, playerStat := range response.Response {
		playerStatEntity := players.PlayerStatisticEntity{}

		err := e.mapPlayerStatistic(playerStat, &playerStatEntity)
		if err != nil {
			return fmt.Errorf("error mapping player statistic for player %d: %w", playerStat.Player.Id, err)
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

func (e *EntityTransformer) mapPlayerStatistic(player PlayerStatsEntity, playerBusinessEntity *players.PlayerStatisticEntity) error {
	secondsPlayed := 0
	if player.Minutes != "" {
		splittedGameTime := strings.Split(player.Minutes, ":")
		minutesPlayed, err := strconv.Atoi(splittedGameTime[0])
		if err == nil {
			secondsAfterMinutes := 0
			if len(splittedGameTime) > 1 {
				secondsAfterMinutes, _ = strconv.Atoi(splittedGameTime[1])
			}
			secondsPlayed = minutesPlayed*60 + secondsAfterMinutes
		}
	}

	fgp := 0.0
	if player.FieldGoals.Percentage != nil {
		fgp = *player.FieldGoals.Percentage
	} else if player.FieldGoals.Attempts != nil && *player.FieldGoals.Attempts > 0 && player.FieldGoals.Total != nil {
		fgp = (float64(*player.FieldGoals.Total) / float64(*player.FieldGoals.Attempts)) * 100
	}

	*playerBusinessEntity = players.PlayerStatisticEntity{
		PlayerExternalId: strconv.Itoa(player.Player.Id),
		PlayerModel: players.PlayerModel{
			FullName: player.Player.Name,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds:        secondsPlayed,
			PlsMin:               0, // api-basketball doesn't seem to have +/- in player stats in the example I got
			Points:               uint8(player.Points),
			Rebounds:             uint8(player.Rebounds.Total),
			Assists:              uint8(player.Assists),
			Steals:               0, // not in basic player stats example
			Blocks:               0, // not in basic player stats example
			FieldGoalsPercentage: statsutil.FromPercentage100(fgp),
			Turnovers:            0, // not in basic player stats example
		},
	}

	return nil
}

func NewEntityTransformer() EntityTransformer {
	return EntityTransformer{}
}
