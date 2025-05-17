package statistics

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/statistics/cdn_nba"
	"IMP/app/log"
	"IMP/app/pkg/array_utils"
	"IMP/app/pkg/time_utils"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type nbaMapper struct {
	league    *domain.League
	timeUtils time_utils.TimeUtilsInterface

	logger log.Logger
}

func newNbaMapper(league *domain.League) *nbaMapper {
	return &nbaMapper{
		league:    league,
		timeUtils: time_utils.NewTimeUtils(),
		logger:    log.NewLogger(),
	}
}

func (c *nbaMapper) mapGame(gameDto cdn_nba.BoxScoreDto) (GameBoxScoreDTO, error) {
	// calculate full game duration
	duration := 0
	duration = c.league.PeriodsNumber * c.league.PeriodDuration
	for i := 0; i < gameDto.Period-c.league.PeriodsNumber; i++ {
		duration += c.league.OvertimeDuration
	}
	homeTeam, err := c.mapTeam(gameDto.HomeTeam)

	if err != nil {
		return GameBoxScoreDTO{}, err
	}
	awayTeam, err := c.mapTeam(gameDto.AwayTeam)

	if err != nil {
		return GameBoxScoreDTO{}, err
	}

	gameBoxScoreDto := GameBoxScoreDTO{
		Id:            gameDto.GameId,
		LeagueAliasEn: c.league.AliasEn,
		IsFinal:       gameDto.GameStatus == 3,
		HomeTeam:      homeTeam,
		AwayTeam:      awayTeam,
		PlayedMinutes: duration,
		ScheduledAt:   gameDto.GameTimeUTC,
	}

	return gameBoxScoreDto, nil
}

func (c *nbaMapper) mapTeam(dto cdn_nba.TeamBoxScoreDto) (TeamBoxScoreDTO, error) {
	dtos, err := array_utils.Map(dto.Players, func(player cdn_nba.PlayerBoxScoreDto) (PlayerDTO, error) {
		return c.mapPlayer(player)
	})
	if err != nil {
		return TeamBoxScoreDTO{}, err
	}

	return TeamBoxScoreDTO{
		Alias:    dto.TeamTricode,
		Name:     dto.TeamName,
		LeagueId: strconv.Itoa(dto.TeamId),
		Scored:   dto.Score,
		Players:  dtos,
	}, nil
}

func (c *nbaMapper) mapPlayer(dto cdn_nba.PlayerBoxScoreDto) (PlayerDTO, error) {
	seconds, err := c.timeUtils.FormattedMinutesToSeconds(dto.Statistics.Minutes, playedTimeFormat)
	if err != nil {
		return PlayerDTO{}, err
	}
	return PlayerDTO{
		FullNameLocal:  dto.Name,
		BirthDate:      nil,
		LeaguePlayerID: strconv.Itoa(dto.PersonId),
		Statistic: PlayerStatisticDTO{
			PlsMin:        int(dto.Statistics.Plus - dto.Statistics.Minus),
			PlayedSeconds: seconds,
			IsBench:       dto.Starter != "1",
		},
	}, nil
}

const playedTimeFormat = "PT%mM%sS"

type nbaProvider struct {
	cdnNbaClient *cdn_nba.Client
	mapper       *nbaMapper
}

func (n *nbaProvider) GamesByDate(date time.Time) ([]string, error) {
	schedule := n.cachedSeasonSchedule()

	formattedSearchedDate := date.Format("01/02/2006 00:00:00")

	for _, gameDate := range schedule.Games {
		if gameDate.GameDate == formattedSearchedDate {
			return array_utils.Map(gameDate.Games, func(game cdn_nba.GameSeasonScheduleDto) (string, error) {
				return game.GameId, nil
			})
		}
	}

	return make([]string, 0), nil
}

func (n *nbaProvider) GameBoxScore(gameId string) (*GameBoxScoreDTO, error) {
	gameDto := n.cdnNbaClient.BoxScore(gameId)

	gameBoxScoreDto, err := n.mapper.mapGame(gameDto)

	return &gameBoxScoreDto, err
}

func (n *nbaProvider) GamesByTeam(teamId string) ([]string, error) {
	panic("implement me")
}

// cachedSeasonSchedule returns the cached season schedule from a file if the file exists and is not older than 7 days
//
// Otherwise it makes a request to the CDN NBA API
func (n *nbaProvider) cachedSeasonSchedule() cdn_nba.SeasonScheduleDto {
	cacheFilePath := filepath.Join(os.TempDir(), "nba_schedule_cache.json")

	// Checks if cached file exists and is not older than 7 days
	if info, err := os.Stat(cacheFilePath); err == nil {
		if time.Since(info.ModTime()) < 7*time.Hour*24 {
			data, err := os.ReadFile(cacheFilePath)
			if err == nil {
				var schedule cdn_nba.SeasonScheduleDto
				if json.Unmarshal(data, &schedule) == nil {
					return schedule
				}
			}
		}
	}

	n.mapper.logger.Info("There is no cached file or it is outdated, making a request...")

	// Making request to get the schedule
	schedule := n.cdnNbaClient.ScheduleSeason()

	data, err := json.Marshal(schedule)
	if err == nil {
		n.mapper.logger.Info("Saving schedule to cache...")
		// Even if there is an error, we still return the schedule from response
		_ = os.WriteFile(cacheFilePath, data, 0644)
	}

	return schedule
}

func newNbaProvider(league *domain.League) *nbaProvider {
	return &nbaProvider{
		cdnNbaClient: cdn_nba.NewCdnNbaClient(),
		mapper:       newNbaMapper(league),
	}
}

// todo: как будто можно абстрагировать мапперы, или что-то такое
