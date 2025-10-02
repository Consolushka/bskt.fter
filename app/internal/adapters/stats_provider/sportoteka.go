package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/sportoteka"
	"strconv"
	"time"
)

type SportotekaStatsProviderAdapter struct {
	tag  string
	year int

	client      sportoteka.ClientInterface
	transformer sportoteka.EntityTransformer
}

func (s SportotekaStatsProviderAdapter) GetGamesStatsByDate(date time.Time) ([]games.GameStatEntity, error) {
	calendar, err := s.client.Calendar(s.tag, s.year, date, date)
	if err != nil {
		return make([]games.GameStatEntity, 0), err
	}

	gamesEntities := make([]games.GameStatEntity, calendar.TotalCount)

	for i, calendarGame := range calendar.Items {
		gameBoxScore, err := s.client.BoxScore(strconv.Itoa(calendarGame.Game.Id))
		if err != nil {
			return []games.GameStatEntity{}, err
		}

		entity, err := s.transformer.Transform(gameBoxScore.Result)
		if err != nil {
			return make([]games.GameStatEntity, 0), err
		}

		gamesEntities[i] = entity
	}

	return gamesEntities, nil
}

func NewSportotekaStatsProvider(client sportoteka.ClientInterface, tag string, year int) SportotekaStatsProviderAdapter {
	return SportotekaStatsProviderAdapter{
		tag:         tag,
		year:        year,
		client:      client,
		transformer: sportoteka.EntityTransformer{},
	}
}
