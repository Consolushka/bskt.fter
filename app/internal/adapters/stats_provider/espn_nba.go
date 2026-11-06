package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/espn"
	"IMP/app/internal/ports"
	"fmt"
	"time"

	// Дни ESPN считаются по нью-йоркскому времени. Описание этого часового пояса встраиваем
	// в саму программу, чтобы не зависеть от того, есть ли база часовых поясов на сервере.
	_ "time/tzdata"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
)

var _ ports.StatsProvider = (*EspnNbaStatsProviderAdapter)(nil)

type EspnNbaStatsProviderAdapter struct {
	client            espn.ClientInterface
	entityTransformer espn.EntityTransformer
	location          *time.Location // America/New_York
}

func NewEspnNbaStatsProviderAdapter(client espn.ClientInterface) (EspnNbaStatsProviderAdapter, error) {
	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return EspnNbaStatsProviderAdapter{}, fmt.Errorf("load America/New_York location: %w", err)
	}

	return EspnNbaStatsProviderAdapter{
		client:            client,
		entityTransformer: espn.NewEntityTransformer(),
		location:          location,
	}, nil
}

func (a EspnNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	start := from.In(a.location)
	end := to.In(a.location)
	day := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, a.location)
	last := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, a.location)

	result := make([]games.GameStatEntity, 0)
	for d := day; !d.After(last); d = d.AddDate(0, 0, 1) {
		date := d.Format("20060102")

		resp, err := a.client.Scoreboard(date)
		if err != nil {
			return nil, fmt.Errorf("scoreboard for %s returned error: %w", date, err)
		}

		for _, event := range resp.Events {
			game, ok, err := a.entityTransformer.TransformEvent(event)
			if err != nil {
				// Одна битая игра не должна останавливать весь турнир.
				composite_logger.Error("Failed to parse ESPN event", map[string]interface{}{
					"eventId": event.Id,
					"date":    event.Date,
					"error":   err,
				})
				continue
			}
			if !ok {
				continue
			}
			result = append(result, game)
		}
	}

	return result, nil
}

func (a EspnNbaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	if game.ExternalGameId == "" {
		return game, nil
	}

	// На любой ошибке возвращаем входную игру: процессор пишет предупреждение с полями из неё.
	resp, err := a.client.Summary(game.ExternalGameId)
	if err != nil {
		return game, fmt.Errorf("ESPN event %s (%s): summary returned error: %w", game.ExternalGameId, game.GameModel.Title, err)
	}

	err = a.entityTransformer.MapBoxScore(resp, &game)
	if err != nil {
		return game, fmt.Errorf("ESPN event %s (%s): transform returned error: %w", game.ExternalGameId, game.GameModel.Title, err)
	}

	err = a.entityTransformer.VerifyTotals(game)
	if err != nil {
		return game, fmt.Errorf("ESPN event %s (%s): verify returned error: %w", game.ExternalGameId, game.GameModel.Title, err)
	}

	return game, nil
}

func (a EspnNbaStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	resp, err := a.client.Athlete(id)
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("ESPN athlete %s: %w", id, err)
	}

	bio, err := a.entityTransformer.MapAthleteBio(resp)
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("ESPN athlete %s: map returned error: %w", id, err)
	}

	return bio, nil
}
