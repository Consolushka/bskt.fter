package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/api_nba"
	"errors"
	"strconv"
	"time"
)

type ApiNbaStatsProviderAdapter struct {
	client api_nba.ClientInterface

	entityTransformer api_nba.EntityTransformer
}

func (a ApiNbaStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	entity := players.PlayerBioEntity{}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return players.PlayerBioEntity{}, err
	}
	playerBio, err := a.client.PlayerInfo(intId, "", 0, 0, "", "")
	if err != nil {
		return players.PlayerBioEntity{}, err
	}

	entity.BirthDate, err = time.Parse("2006-01-02", playerBio.Response[0].Birth.Date)
	if err != nil {
		entity.BirthDate = time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
		return entity, errors.New(err.Error())
	}

	return entity, nil
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
		gameStatEntity := a.entityTransformer.TransformWithoutPlayers(passedGame)
		gameStatEntity.ExternalGameId = strconv.Itoa(passedGame.Id)
		gameStatEntity.HomeTeamExternalId = passedGame.Teams.Home.Id
		gameStatEntity.AwayTeamExternalId = passedGame.Teams.Visitors.Id

		gamesStatsEntities = append(gamesStatsEntities, gameStatEntity)
	}

	return gamesStatsEntities, nil
}

func (a ApiNbaStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	if game.ExternalGameId == "" {
		return game, nil
	}

	gameId, err := strconv.Atoi(game.ExternalGameId)
	if err != nil {
		return games.GameStatEntity{}, err
	}

	err = a.entityTransformer.EnrichGamePlayers(gameId, game.HomeTeamExternalId, game.AwayTeamExternalId, &game)
	if err != nil {
		return games.GameStatEntity{}, err
	}

	return game, nil
}

func NewApiNbaStatsProviderAdapter(client api_nba.ClientInterface) ApiNbaStatsProviderAdapter {
	return ApiNbaStatsProviderAdapter{
		client:            client,
		entityTransformer: api_nba.NewEntityTransformer(client),
	}
}
