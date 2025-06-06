package statistics

import (
	"IMP/app/internal/domain"
	"IMP/app/internal/statistics/infobasket"
	"IMP/app/pkg/array_utils"
	"IMP/app/pkg/string_utils"
	"IMP/app/pkg/translator"
	"errors"
	"strconv"
	"time"
)

type mlblMapperInterface interface {
	mapGame(game infobasket.GameBoxScoreResponse) (*GameBoxScoreDTO, error)
	mapTeam(teamBoxScore infobasket.TeamBoxScoreDto) (TeamBoxScoreDTO, error)
	mapPlayer(player infobasket.PlayerBoxScoreDto) (PlayerDTO, error)
}

type mlblMapper struct {
	league      *domain.League
	stringUtils string_utils.StringUtilsInterface
	translator  translator.Interface
}

func newMlblMapper(utilsInterface string_utils.StringUtilsInterface, translator translator.Interface, league *domain.League) *mlblMapper {
	return &mlblMapper{
		league:      league,
		stringUtils: utilsInterface,
		translator:  translator,
	}
}

func (m *mlblMapper) mapGame(game infobasket.GameBoxScoreResponse) (*GameBoxScoreDTO, error) {
	duration := 0
	duration = m.league.PeriodsNumber * m.league.PeriodDuration
	for i := 0; i < game.MaxPeriod-m.league.PeriodsNumber; i++ {
		duration += m.league.OvertimeDuration
	}

	scheduled, err := time.Parse("02.01.2006 15.04", game.GameDate+" "+game.GameTime)
	if err != nil {
		return nil, errors.New("can't parse game datetime. given game datetime: " + game.GameDate + " " + game.GameTime + " doesn't match format 02.01.2006 15.04")
	}

	homeTeamDto, err := m.mapTeam(game.GameTeams[0])
	if err != nil {
		return nil, err
	}
	awayTeamDto, err := m.mapTeam(game.GameTeams[1])
	if err != nil {
		return nil, err
	}

	gameBoxScoreDto := GameBoxScoreDTO{
		LeagueAliasEn: m.league.AliasEn,
		IsFinal:       game.GameStatus == 1,
		HomeTeam:      homeTeamDto,
		AwayTeam:      awayTeamDto,
		PlayedMinutes: duration,
		ScheduledAt:   scheduled,
	}

	return &gameBoxScoreDto, nil
}

func (m *mlblMapper) mapTeam(teamBoxScore infobasket.TeamBoxScoreDto) (TeamBoxScoreDTO, error) {
	playersDtos, err := array_utils.Map(teamBoxScore.Players, func(player infobasket.PlayerBoxScoreDto) (PlayerDTO, error) {
		return m.mapPlayer(player)
	})
	if err != nil {
		return TeamBoxScoreDTO{}, err
	}

	return TeamBoxScoreDTO{
		Alias:    teamBoxScore.TeamName.CompTeamAbcNameEn,
		Name:     teamBoxScore.TeamName.CompTeamNameEn,
		LeagueId: strconv.Itoa(teamBoxScore.TeamID),
		Scored:   teamBoxScore.Score,
		Players:  playersDtos,
	}, nil
}

func (m *mlblMapper) mapPlayer(player infobasket.PlayerBoxScoreDto) (PlayerDTO, error) {
	birthdate, err := time.Parse("02.01.2006", player.PersonBirth)
	if err != nil {
		return PlayerDTO{}, errors.New("can't parse player birthdate. given birthdate: " + player.PersonBirth + " doesn't match format 02.01.2006")
	}

	var enPersonName string

	hasNonLatinChars, err := m.stringUtils.HasNonLanguageChars(player.PersonNameEn, string_utils.Latin)
	if err != nil {
		return PlayerDTO{}, err
	}

	if hasNonLatinChars {
		ruCode := "ru"
		enPersonName = m.translator.Translate(player.PersonNameEn, &ruCode, "en")
	} else {
		enPersonName = player.PersonNameEn
	}

	return PlayerDTO{
		FullNameLocal:  player.PersonNameRu,
		FullNameEn:     enPersonName,
		BirthDate:      &birthdate,
		LeaguePlayerID: strconv.Itoa(player.PersonID),
		Statistic: PlayerStatisticDTO{
			PlsMin:        player.PlusMinus,
			PlayedSeconds: player.Seconds,
			IsBench:       !player.IsStart,
		},
	}, nil
}

type mlblProvider struct {
	client infobasket.ClientInterface
	mapper mlblMapperInterface
}

func (i *mlblProvider) GameBoxScore(gameId string) (*GameBoxScoreDTO, error) {
	gameDto := i.client.BoxScore(gameId)

	game, err := i.mapper.mapGame(gameDto)
	if err != nil {
		return nil, err
	}

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

	gamesIds, err := array_utils.Map(scheduleJson.Games, func(game infobasket.GameScheduleDto) (string, error) {
		if game.GameStatus == 1 {
			return strconv.Itoa(game.GameID), nil
		}

		return "", errors.New("game is not final. or game status is: " + strconv.Itoa(game.GameStatus))
	})
	if err != nil {
		return nil, err
	}

	return array_utils.Filter(gamesIds, func(gameId string) bool {
		return gameId != ""
	}), nil
}

func newMlblProvider(league *domain.League) *mlblProvider {
	return &mlblProvider{
		client: infobasket.NewInfobasketClient(),
		mapper: newMlblMapper(string_utils.NewStringUtils(), translator.NewTranslator(), league),
	}
}
