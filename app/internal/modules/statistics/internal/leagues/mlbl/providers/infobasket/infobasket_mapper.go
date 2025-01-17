package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	"IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"time"
)

type mapper struct{}

func newMapper() *mapper {
	return &mapper{}
}

func (m *mapper) mapGame(game boxscore.GameInfo) *models.GameBoxScoreDTO {
	league := enums.MLBL

	duration := 0
	duration = 4 * league.QuarterDuration()
	for i := 0; i < game.MaxPeriod-4; i++ {
		duration += league.OvertimeDuration()
	}

	scheduled, _ := time.Parse("2006-01-02 23.1", game.GameDate+" "+game.GameTime)

	gameBoxScoreDto := models.GameBoxScoreDTO{
		League:        league,
		HomeTeam:      m.mapTeam(game.GameTeams[0]),
		AwayTeam:      m.mapTeam(game.GameTeams[1]),
		PlayedMinutes: game.MaxPeriod,
		ScheduledAt:   scheduled,
	}

	return &gameBoxScoreDto
}

func (m *mapper) mapTeam(teamBoxscore boxscore.TeamBoxscore) models.TeamBoxScoreDTO {
	return models.TeamBoxScoreDTO{
		Alias:    teamBoxscore.TeamName.CompTeamAbcNameEn,
		Name:     teamBoxscore.TeamName.CompTeamNameEn,
		LeagueId: teamBoxscore.TeamID,
		Scored:   teamBoxscore.Score,
		Players: array_utils.Map(teamBoxscore.Players, func(player boxscore.PlayerBoxscore) models.PlayerDTO {
			return m.mapPlayer(player)
		}),
	}
}

func (m *mapper) mapPlayer(player boxscore.PlayerBoxscore) models.PlayerDTO {
	birthdate, _ := time.Parse("2006-01-02", player.PersonBirth)

	return models.PlayerDTO{
		FullName:       player.PersonNameRu,
		BirthDate:      &birthdate,
		LeaguePlayerID: player.PersonID,
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        player.PlusMinus,
			PlayedSeconds: player.Seconds,
			IsBench:       !player.IsStart,
		},
	}
}
