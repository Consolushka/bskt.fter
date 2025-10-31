package sportoteka

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/logger"
	"strconv"
	"time"
)

type EntityTransformer struct {
}

func (e *EntityTransformer) Transform(game GameBoxScoreEntity) (games.GameStatEntity, error) {
	var homeTeamStats, awayTeamStats TeamBoxScoreEntity
	for _, team := range game.Teams {
		if team.TeamNumber == 1 {
			homeTeamStats = team
		}

		if team.TeamNumber == 2 {
			awayTeamStats = team
		}
	}

	return games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: game.Game.ScheduledTime,
			Title:       game.Team1.AbcName + " - " + game.Team2.AbcName,
		},
		HomeTeamStat: e.teamTransform(game.Team1, homeTeamStats, awayTeamStats.Total.Points),
		AwayTeamStat: e.teamTransform(game.Team2, awayTeamStats, homeTeamStats.Total.Points),
	}, nil
}

func (e *EntityTransformer) teamTransform(teamInfo TeamInfoEntity, teamBoxScore TeamBoxScoreEntity, opponentsScore int) teams.TeamStatEntity {
	playerStats := make([]players.PlayerStatisticEntity, len(teamBoxScore.Starts))

	for i, player := range teamBoxScore.Starts {
		if player.StartRole == "Team" {
			continue
		}

		playerStat, err := e.playerTransform(player)
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
			Name:     teamInfo.Name,
			HomeTown: teamInfo.RegionName,
		},
		GameTeamStatModel: teams.GameTeamStatModel{
			Score:     teamBoxScore.Total.Points,
			FinalDiff: teamBoxScore.Total.Points - opponentsScore,
		},
		PlayerStats: playerStats,
	}
}

func (e *EntityTransformer) playerTransform(player TeamBoxScoreStartEntity) (players.PlayerStatisticEntity, error) {
	parsedBirth, err := time.Parse(time.RFC3339, player.Birthday+"+03:00")
	if err != nil {
		return players.PlayerStatisticEntity{}, err
	}

	return players.PlayerStatisticEntity{
		PlayerExternalId: strconv.Itoa(*player.PersonId),
		PlayerModel: players.PlayerModel{
			FullName:  player.LastName + " " + player.FirstName,
			BirthDate: parsedBirth,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds: player.Stats.Second,
			PlsMin:        int8(player.Stats.PlusMinus),
		},
	}, nil
}
