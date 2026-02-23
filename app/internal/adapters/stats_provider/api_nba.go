package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/api_nba"
	"errors"
	"fmt"
	"reflect"
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
		return players.PlayerBioEntity{}, fmt.Errorf("atoi with %s returned error: %w", id, err)
	}
	playerBio, err := a.client.PlayerInfo(intId, "", 0, 0, "", "")
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("playerInfo with %v, %v, %v, %v, %v, %v from %s returned error: %w", intId, "", 0, 0, "", "", reflect.TypeOf(a.client), err)
	}

	if len(playerBio.Response) == 0 {
		return players.PlayerBioEntity{}, errors.New("empty player info response")
	}

	entity.FullName = playerBio.Response[0].Firstname + " " + playerBio.Response[0].Lastname
	entity.BirthDate, err = time.Parse("2006-01-02", playerBio.Response[0].Birth.Date)
	if err != nil {
		entity.BirthDate = time.Date(1, 1, 1, 1, 1, 1, 1, time.UTC)
		return entity, fmt.Errorf("time.Parse with %s, %v returned error: %w", "2006-01-02", playerBio.Response[0].Birth.Date, err)
	}

	return entity, nil
}

func (a ApiNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	var passedGames []api_nba.GameEntity

	if from.Truncate(24*time.Hour) != to.Truncate(24*time.Hour) {
		fromResponse, err := a.client.Games(0, from.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, fmt.Errorf("games with %v, %v, %v, %v, %v, %v from %s returned error: %w", 0, from.Format("2006-01-02"), "1", "", "", "", reflect.TypeOf(a.client), err)
		}
		passedGames = fromResponse.Response

		toResponse, err := a.client.Games(0, to.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, fmt.Errorf("games with %v, %v, %v, %v, %v, %v from %s returned error: %w", 0, to.Format("2006-01-02"), "1", "", "", "", reflect.TypeOf(a.client), err)
		}
		passedGames = append(passedGames, toResponse.Response...)
	} else {
		response, err := a.client.Games(0, to.Format("2006-01-02"), "1", "", "", "")
		if err != nil {
			return nil, fmt.Errorf("games with %v, %v, %v, %v, %v, %v from %s returned error: %w", 0, to.Format("2006-01-02"), "1", "", "", "", reflect.TypeOf(a.client), err)
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
		return games.GameStatEntity{}, fmt.Errorf("atoi with %v returned error: %w", game.ExternalGameId, err)
	}

	playerStatsResponse, err := a.client.PlayersStatistics(0, gameId, 0, "")
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("PlayersStatistics with %v from %s returned error: %w", gameId, reflect.TypeOf(a.client), err)
	}

	err = a.entityTransformer.MapPlayerStatistics(playerStatsResponse, game.HomeTeamExternalId, game.AwayTeamExternalId, &game)
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("MapPlayerStatistics returned error: %w", err)
	}

	return game, nil
}

func NewApiNbaStatsProviderAdapter(client api_nba.ClientInterface) ApiNbaStatsProviderAdapter {
	return ApiNbaStatsProviderAdapter{
		client:            client,
		entityTransformer: api_nba.NewEntityTransformer(),
	}
}
