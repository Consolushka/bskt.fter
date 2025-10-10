package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/infra/api_nba"
	"time"
)

type ApiNbaStatsProviderAdapter struct {
	client api_nba.ClientInterface

	entityTransformer api_nba.EntityTransformer
}

func (a ApiNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	var passedGames []api_nba.GameEntity

	if from.Truncate(24*time.Hour) != to.Truncate(24*time.Hour) {
		fromResponse, err := a.client.Games(0, from.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, err
		}
		passedGames = fromResponse.Response

		toResponse, err := a.client.Games(0, to.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, err
		}
		passedGames = append(passedGames, toResponse.Response...)
	} else {
		response, err := a.client.Games(0, to.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, err
		}

		passedGames = response.Response
	}

	gamesStatsEntities := make([]games.GameStatEntity, 0, len(passedGames))

	for _, passedGame := range passedGames {
		if passedGame.Status.Short != 3 {
			continue
		}
		transformer := api_nba.NewEntityTransformer(a.client)

		gameStatEntity, err := transformer.Transform(passedGame)
		if err != nil {
			return nil, err
		}

		gamesStatsEntities = append(gamesStatsEntities, gameStatEntity)
	}

	return gamesStatsEntities, nil
}

func NewApiNbaStatsProviderAdapter(client api_nba.ClientInterface) ApiNbaStatsProviderAdapter {
	return ApiNbaStatsProviderAdapter{
		client:            client,
		entityTransformer: api_nba.NewEntityTransformer(client),
	}
}
