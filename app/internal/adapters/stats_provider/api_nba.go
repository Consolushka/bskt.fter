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

func (a ApiNbaStatsProviderAdapter) GetGamesStatsByDate(date time.Time) ([]games.GameStatEntity, error) {
	passedGames, err := a.client.Games(0, date.Format("2006-01-02"), "1", "", "", "")
	if err != nil {
		return nil, err
	}

	gamesStatsEntities := make([]games.GameStatEntity, 0, len(passedGames.Response))

	for _, passedGame := range passedGames.Response {
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
