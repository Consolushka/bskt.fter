package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/sportoteka"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type SportotekaStatsProviderAdapter struct {
	tag  string
	year int

	client      sportoteka.ClientInterface
	transformer sportoteka.EntityTransformer
}

func (s SportotekaStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	//TODO implement me
	panic("implement me")
}

func (s SportotekaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	calendar, err := s.client.Calendar(s.tag, s.year, from, to)
	if err != nil {
		return make([]games.GameStatEntity, 0), err
	}

	gamesEntities := make([]games.GameStatEntity, 0, calendar.TotalCount)

	for _, calendarGame := range calendar.Items {
		if calendarGame.Game.GameStatus != "Result" {
			continue
		}
		gameBoxScore, err := s.client.BoxScore(strconv.Itoa(calendarGame.Game.Id))
		if err != nil {
			return []games.GameStatEntity{}, fmt.Errorf("BoxScore with %v from %s returned error: %w", strconv.Itoa(calendarGame.Game.Id), reflect.TypeOf(s.client), err)
		}

		if gameBoxScore.Result.Game.GameStatus != "ResultConfirmed" && gameBoxScore.Result.Game.GameStatus != "Result" {
			continue
		}

		entity, err := s.transformer.Transform(gameBoxScore.Result)
		if err != nil {
			return make([]games.GameStatEntity, 0), err
		}

		entity.GameModel.Title = calendarGame.Team1.AbcName + " - " + calendarGame.Team2.AbcName
		entity.HomeTeamStat.TeamModel.Name = calendarGame.Team1.Name
		entity.AwayTeamStat.TeamModel.Name = calendarGame.Team2.Name

		gamesEntities = append(gamesEntities, entity)
	}

	return gamesEntities, nil
}

func (s SportotekaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	return game, nil
}

func NewSportotekaStatsProvider(client sportoteka.ClientInterface, tag string, year int) SportotekaStatsProviderAdapter {
	return SportotekaStatsProviderAdapter{
		tag:         tag,
		year:        year,
		client:      client,
		transformer: sportoteka.EntityTransformer{},
	}
}
