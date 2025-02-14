package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba/dtos/boxscore"
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"IMP/app/log"
	"strconv"
)

type mapper struct {
	leagueRepository *leaguesDomain.Repository
}

func newMapper() *mapper {
	return &mapper{
		leagueRepository: leaguesDomain.NewRepository(),
	}
}

func (c *mapper) mapGame(gameDto boxscore.GameDTO) models.GameBoxScoreDTO {
	league, err := c.leagueRepository.GetLeagueByAliasEn("nba")
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}
	// calculate full game duration
	duration := 0
	duration = league.PeriodsNumber * league.PeriodDuration
	for i := 0; i < gameDto.Period-league.PeriodsNumber; i++ {
		duration += league.OvertimeDuration
	}
	gameBoxScoreDto := models.GameBoxScoreDTO{
		Id:            gameDto.GameId,
		LeagueAliasEn: league.AliasEn,
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
		LeagueId: strconv.Itoa(dto.TeamId),
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
		LeaguePlayerID: strconv.Itoa(dto.PersonId),
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        dto.Statistics.Plus - dto.Statistics.Minus,
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(dto.Statistics.Minutes, playedTimeFormat),
			IsBench:       dto.Starter != "1",
		},
	}
}
