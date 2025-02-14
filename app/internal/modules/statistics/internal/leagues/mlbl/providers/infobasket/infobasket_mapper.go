package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket/dtos/boxscore"
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/log"
	"strconv"
	"time"
)

type mapper struct {
	leagueRepository *leaguesDomain.Repository
}

func newMapper() *mapper {
	return &mapper{
		leagueRepository: leaguesDomain.NewRepository(),
	}
}

func (m *mapper) mapGame(game boxscore.GameInfo) *models.GameBoxScoreDTO {
	league, err := m.leagueRepository.GetLeagueByAliasEn("MLBL")
	if err != nil {
		log.Fatalln(err)
		panic(err)
	}

	duration := 0
	duration = league.PeriodsNumber * league.PeriodDuration
	for i := 0; i < game.MaxPeriod-league.PeriodsNumber; i++ {
		duration += league.OvertimeDuration
	}

	scheduled, _ := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTime)

	gameBoxScoreDto := models.GameBoxScoreDTO{
		LeagueAliasEn: league.AliasEn,
		HomeTeam:      m.mapTeam(game.GameTeams[0]),
		AwayTeam:      m.mapTeam(game.GameTeams[1]),
		PlayedMinutes: duration,
		ScheduledAt:   scheduled,
	}

	return &gameBoxScoreDto
}

func (m *mapper) mapTeam(teamBoxScore boxscore.TeamBoxscore) models.TeamBoxScoreDTO {
	return models.TeamBoxScoreDTO{
		Alias:    teamBoxScore.TeamName.CompTeamAbcNameEn,
		Name:     teamBoxScore.TeamName.CompTeamNameEn,
		LeagueId: strconv.Itoa(teamBoxScore.TeamID),
		Scored:   teamBoxScore.Score,
		Players: array_utils.Map(teamBoxScore.Players, func(player boxscore.PlayerBoxscore) models.PlayerDTO {
			return m.mapPlayer(player)
		}),
	}
}

func (m *mapper) mapPlayer(player boxscore.PlayerBoxscore) models.PlayerDTO {
	birthdate, _ := time.Parse("02.01.2006", player.PersonBirth)

	return models.PlayerDTO{
		FullNameLocal:  player.PersonNameRu,
		FullNameEn:     player.PersonNameEn,
		BirthDate:      &birthdate,
		LeaguePlayerID: strconv.Itoa(player.PersonID),
		Statistic: models.PlayerStatisticDTO{
			PlsMin:        player.PlusMinus,
			PlayedSeconds: player.Seconds,
			IsBench:       !player.IsStart,
		},
	}
}
