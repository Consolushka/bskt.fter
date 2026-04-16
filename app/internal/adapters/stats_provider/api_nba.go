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
		return entity, fmt.Errorf("time.Parse with %s, %v for api_nba returned error: %w", "2006-01-02", playerBio.Response[0].Birth.Date, err)
	}

	return entity, nil
}

func (a ApiNbaStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	var rawGames []api_nba.GameEntity

	// Calculate all unique dates to fetch.
	uniqueDates := make(map[string]struct{})

	// Iterate through all dates in the range [from, to]
	for d := from; !d.After(to); d = d.Add(24 * time.Hour) {
		uniqueDates[d.Format("2006-01-02")] = struct{}{}
	}
	// Explicitly add 'to' date to ensure it's covered even if loop boundary is tricky
	uniqueDates[to.Format("2006-01-02")] = struct{}{}

	for dateStr := range uniqueDates {
		response, err := a.client.Games(0, dateStr, "", "", "", "")
		if err != nil {
			return nil, fmt.Errorf("games with date %s from %s returned error: %w", dateStr, reflect.TypeOf(a.client), err)
		}
		rawGames = append(rawGames, response.Response...)
	}

	gamesStatsEntities := make([]games.GameStatEntity, 0, len(rawGames))
	processedGameIds := make(map[int]struct{})

	for _, rawGame := range rawGames {
		// Deduplicate: same game might be returned for different dates if we overlap
		if _, ok := processedGameIds[rawGame.Id]; ok {
			continue
		}
		processedGameIds[rawGame.Id] = struct{}{}

		if rawGame.Status.Short != 3 {
			continue
		}
		gameStatEntity := a.entityTransformer.TransformWithoutPlayers(rawGame)
		gameStatEntity.ExternalGameId = strconv.Itoa(rawGame.Id)
		gameStatEntity.HomeTeamExternalId = rawGame.Teams.Home.Id
		gameStatEntity.AwayTeamExternalId = rawGame.Teams.Visitors.Id

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
