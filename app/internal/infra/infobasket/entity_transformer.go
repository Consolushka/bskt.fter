package infobasket

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/logger"
	"time"
)

type EntityTransformer struct{}

func (e *EntityTransformer) Transform(game GameBoxScoreResponse) (games.GameStatEntity, error) {
	parse, err := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTimeMsk)
	if err != nil {
		return games.GameStatEntity{}, err
	}

	return games.GameStatEntity{
		GameModel: games.GameModel{
			ScheduledAt: parse,
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
		return players.PlayerStatisticEntity{}, err
	}

	playerName := playerStat.PersonNameEn
	if playerName == "New Player" {
		playerName = playerStat.PersonNameRu
	}

	return players.PlayerStatisticEntity{
		PlayerModel: players.PlayerModel{
			FullName:  playerName,
			BirthDate: parsedBirth,
		},
		GameTeamPlayerStatModel: players.GameTeamPlayerStatModel{
			PlayedSeconds: playerStat.Seconds,
			PlsMin:        int8(playerStat.PlusMinus),
		},
	}, nil
}
