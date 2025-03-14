package mlbl

import (
	"IMP/app/internal/infrastructure/infobasket"
	"IMP/app/internal/modules/statistics/models"
	"IMP/app/internal/utils/array_utils"
	"strconv"
	"time"
)

type mapper struct{}

func newMapper() *mapper {
	return &mapper{}
}

func (m *mapper) mapGame(game infobasket.GameBoxScoreResponse, regulationPeriodsNumber int, periodDuration int, overtimeDuration int, leagueAlias string) *models.GameBoxScoreDTO {
	duration := 0
	duration = regulationPeriodsNumber * periodDuration
	for i := 0; i < game.MaxPeriod-regulationPeriodsNumber; i++ {
		duration += overtimeDuration
	}

	scheduled, _ := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTime)

	gameBoxScoreDto := models.GameBoxScoreDTO{
		LeagueAliasEn: leagueAlias,
		IsFinal:       game.GameStatus == 1,
		HomeTeam:      m.mapTeam(game.GameTeams[0]),
		AwayTeam:      m.mapTeam(game.GameTeams[1]),
		PlayedMinutes: duration,
		ScheduledAt:   scheduled,
	}

	return &gameBoxScoreDto
}

func (m *mapper) mapTeam(teamBoxScore infobasket.TeamBoxScoreDto) models.TeamBoxScoreDTO {
	return models.TeamBoxScoreDTO{
		Alias:    teamBoxScore.TeamName.CompTeamAbcNameEn,
		Name:     teamBoxScore.TeamName.CompTeamNameEn,
		LeagueId: strconv.Itoa(teamBoxScore.TeamID),
		Scored:   teamBoxScore.Score,
		Players: array_utils.Map(teamBoxScore.Players, func(player infobasket.PlayerBoxScoreDto) models.PlayerDTO {
			return m.mapPlayer(player)
		}),
	}
}

func (m *mapper) mapPlayer(player infobasket.PlayerBoxScoreDto) models.PlayerDTO {
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
