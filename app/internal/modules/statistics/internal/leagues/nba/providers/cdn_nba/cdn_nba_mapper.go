package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
)

type mapper struct{}

func newMapper() *mapper {
	return &mapper{}
}

func (c *mapper) mapGame(gameDto boxscore.GameDTO) models.GameBoxScoreDTO {
	// calculate full game duration
	duration := 0
	duration = 4 * league.QuarterDuration()
	for i := 5; i < gameDto.Period; i++ {
		duration += league.OvertimeDuration()
	}
	gameBoxScoreDto := models.GameBoxScoreDTO{
		League:        league,
		HomeTeam:      c.mapTeam(gameDto.HomeTeam),
		AwayTeam:      c.mapTeam(gameDto.AwayTeam),
		PlayedMinutes: duration,
		ScheduledAt:   gameDto.GameTimeUTC,
	}

	return gameBoxScoreDto
}

func (c *mapper) mapTeam(dto boxscore.TeamDTO) models.TeamBoxScoreDTO {
	return models.TeamBoxScoreDTO{
		Alias:    dto.TeamTricode,
		Name:     dto.TeamName,
		LeagueId: dto.TeamId,
		Scored:   dto.Score,
		Players: array_utils.Map(dto.Players, func(player boxscore.PlayerDTO) models.PlayerDTO {
			return c.mapPlayer(player)
		}),
	}
}

func (c *mapper) mapPlayer(dto boxscore.PlayerDTO) models.PlayerDTO {
	return models.PlayerDTO{
		FullNameLocal:  dto.Name,
		BirthDate:      nil,
		LeaguePlayerID: dto.PersonId,
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        dto.Statistics.Plus - dto.Statistics.Minus,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(dto.Statistics.Minutes, playedTimeFormat),
			IsBench:       dto.Starter != "1",
		},
	}
}
