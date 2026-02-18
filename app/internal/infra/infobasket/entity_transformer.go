package infobasket

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

type EntityTransformer struct{}

func (e *EntityTransformer) Transform(game GameBoxScoreResponse) (games.GameStatEntity, error) {
	parse, err := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTimeMsk)
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("time.Parse with %s, %v,  returned error: %w", "02.01.2006 15.04", game.GameDate+" "+game.GameTimeMsk, err)
	}

	duration := 0
	for i := 0; i < game.MaxPeriod; i++ {
		if i <= 4 {
			duration += 10
		} else {
			duration += 5
		}
	}

	return games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: parse,
			Duration:    duration,
			Title:       game.GameTeams[0].TeamName.CompTeamAbcNameEn + " - " + game.GameTeams[1].TeamName.CompTeamAbcNameEn,
		},
		HomeTeamStat: e.transformTeam(game.GameTeams[0], game.GameTeams[1].Score),
		AwayTeamStat: e.transformTeam(game.GameTeams[1], game.GameTeams[0].Score),
	}, nil
}

func (e *EntityTransformer) transformTeam(team TeamBoxScoreDto, opponentScore int) teams.TeamStatEntity {
	playerStats := make([]players.PlayerStatisticEntity, len(team.Players))

	for i, player := range team.Players {
		playerStat, err := playersTrans(player)
		if err != nil {
			logger.Warn("There was an error with player statistics", map[string]interface{}{
				"player": player,
				"error":  err,
			})
			continue
		}

		playerStats[i] = playerStat
	}

	return teams.TeamStatEntity{
		TeamModel: teams.TeamModel{
			Name:     team.TeamName.CompTeamNameEn,
			HomeTown: team.TeamName.CompTeamRegionNameEn,
		},
		GameTeamStatModel: teams.GameTeamStatModel{
			Score:     team.Score,
			FinalDiff: team.Score - opponentScore,
		},
		PlayerStats: playerStats,
	}
}

func playersTrans(playerStat PlayerBoxScoreDto) (players.PlayerStatisticEntity, error) {
	parsedBirth, err := time.Parse("02.01.2006", playerStat.PersonBirth)
	if err != nil {
		return players.PlayerStatisticEntity{}, fmt.Errorf("time.Parse with %s, %v,  returned error: %w", "02.01.2006", playerStat.PersonBirth, err)
	}

	playerName := playerStat.PersonNameEn
	if playerName == "New Player" {
		playerName = playerStat.PersonNameRu
	}

	playerAttempts := playerStat.Shot2 + playerStat.Shot3
	var percentage float32
	if playerAttempts == 0 {
		percentage = 0
	} else {
		percentage = float32(playerStat.Goal2+playerStat.Goal3) / float32(playerAttempts)
	}
	return players.PlayerStatisticEntity{
		PlayerExternalId: strconv.Itoa(playerStat.PersonID),
		PlayerModel: players.PlayerModel{
			FullName:  playerName,
			BirthDate: parsedBirth,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds:        playerStat.Seconds,
			PlsMin:               int8(playerStat.PlusMinus),
			Points:               uint8(playerStat.Points),
			Rebounds:             uint8(playerStat.Rebound),
			Assists:              uint8(playerStat.Assist),
			Steals:               uint8(playerStat.Steal),
			Blocks:               uint8(playerStat.Blocks),
			Turnovers:            uint8(playerStat.Turnover),
			FieldGoalsPercentage: percentage,
		},
	}, nil
}
