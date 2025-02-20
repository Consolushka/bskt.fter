package cdn_nba

import (
	"IMP/app/internal/infrastructure/cdn_nba"
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"IMP/app/log"
	"strconv"
	"strings"
)

type mapper struct {
	leagueRepository *leaguesDomain.Repository
}

func newMapper() *mapper {
	return &mapper{
		leagueRepository: leaguesDomain.NewRepository(),
	}
}

func (c *mapper) mapGame(gameDto cdn_nba.BoxScoreDto) models.GameBoxScoreDTO {
	league, err := c.leagueRepository.FirstByAliasEn(strings.ToUpper(leaguesModels.NBAAlias))
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

func (c *mapper) mapTeam(dto cdn_nba.TeamBoxScoreDto) models.TeamBoxScoreDTO {
	return models.TeamBoxScoreDTO{
		Alias:    dto.TeamTricode,
		Name:     dto.TeamName,
		LeagueId: strconv.Itoa(dto.TeamId),
		Scored:   dto.Score,
		Players: array_utils.Map(dto.Players, func(player cdn_nba.PlayerBoxScoreDto) models.PlayerDTO {
			return c.mapPlayer(player)
		}),
	}
}

func (c *mapper) mapPlayer(dto cdn_nba.PlayerBoxScoreDto) models.PlayerDTO {
	return models.PlayerDTO{
		FullNameLocal:  dto.Name,
		BirthDate:      nil,
		LeaguePlayerID: strconv.Itoa(dto.PersonId),
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        int(dto.Statistics.Plus - dto.Statistics.Minus),
			PlayedSeconds: time_utils.FormattedMinutesToSeconds(dto.Statistics.Minutes, playedTimeFormat),
			IsBench:       dto.Starter != "1",
		},
	}
}
