package statistics

import (
	"IMP/app/internal/statistics/infobasket"
	"IMP/app/pkg/array_utils"
	"strconv"
	"time"
)

type mlblMapper struct{}

func newMlblMapper() *mlblMapper {
	return &mlblMapper{}
}

func (m *mlblMapper) mapGame(game infobasket.GameBoxScoreResponse, regulationPeriodsNumber int, periodDuration int, overtimeDuration int, leagueAlias string) *GameBoxScoreDTO {
	duration := 0
	duration = regulationPeriodsNumber * periodDuration
	for i := 0; i < game.MaxPeriod-regulationPeriodsNumber; i++ {
		duration += overtimeDuration
	}

	scheduled, _ := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTime)

	gameBoxScoreDto := GameBoxScoreDTO{
		LeagueAliasEn: leagueAlias,
		IsFinal:       game.GameStatus == 1,
		HomeTeam:      m.mapTeam(game.GameTeams[0]),
		AwayTeam:      m.mapTeam(game.GameTeams[1]),
		PlayedMinutes: duration,
		ScheduledAt:   scheduled,
	}

	return &gameBoxScoreDto
}

func (m *mlblMapper) mapTeam(teamBoxScore infobasket.TeamBoxScoreDto) TeamBoxScoreDTO {
	return TeamBoxScoreDTO{
		Alias:    teamBoxScore.TeamName.CompTeamAbcNameEn,
		Name:     teamBoxScore.TeamName.CompTeamNameEn,
		LeagueId: strconv.Itoa(teamBoxScore.TeamID),
		Scored:   teamBoxScore.Score,
		Players: array_utils.Map(teamBoxScore.Players, func(player infobasket.PlayerBoxScoreDto) PlayerDTO {
			return m.mapPlayer(player)
		}),
	}
}

func (m *mlblMapper) mapPlayer(player infobasket.PlayerBoxScoreDto) PlayerDTO {
	birthdate, _ := time.Parse("02.01.2006", player.PersonBirth)

	return PlayerDTO{
		FullNameLocal:  player.PersonNameRu,
		FullNameEn:     player.PersonNameEn,
		BirthDate:      &birthdate,
		LeaguePlayerID: strconv.Itoa(player.PersonID),
		Statistic: PlayerStatisticDTO{
			PlsMin:        player.PlusMinus,
			PlayedSeconds: player.Seconds,
			IsBench:       !player.IsStart,
		},
	}
}

type mlblProvider struct {
	client infobasket.ClientInterface
	mapper *mlblMapper
}

func (i *mlblProvider) GameBoxScore(gameId string) (*GameBoxScoreDTO, error) {
	gameDto := i.client.BoxScore(gameId)

	game := i.mapper.mapGame(gameDto, 4, 10, 5, "MLBL")
	game.Id = gameId
	return game, nil
}

func (i *mlblProvider) GamesByDate(date time.Time) ([]string, error) {
	var result []string
	compIds := []int{89960, 89962}

	for _, compId := range compIds {
		seasonGames := i.client.ScheduledGames(compId)

		for _, game := range seasonGames {
			if game.GameDate == date.Format("02.01.2006") {
				result = append(result, strconv.Itoa(game.GameID))
			}
		}
	}

	return result, nil
}

func (i *mlblProvider) GamesByTeam(teamId string) ([]string, error) {
	scheduleJson := i.client.TeamGames(teamId)

	gamesIds := array_utils.Map(scheduleJson.Games, func(game infobasket.GameScheduleDto) string {
		if game.GameStatus == 1 {
			return strconv.Itoa(game.GameID)
		}

		return ""
	})

	return array_utils.Filter(gamesIds, func(gameId string) bool {
		return gameId != ""
	}), nil
}

func newMlblProvider() *mlblProvider {
	return &mlblProvider{
		client: infobasket.NewInfobasketClient(),
		mapper: newMlblMapper(),
	}
}
