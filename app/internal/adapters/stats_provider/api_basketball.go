package stats_provider

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/infra/api_basketball"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type ApiBasketballStatsProviderAdapter struct {
	client api_basketball.ClientInterface

	entityTransformer api_basketball.EntityTransformer
	leagueId          int
}

func (a ApiBasketballStatsProviderAdapter) GetPlayerBio(id string) (players.PlayerBioEntity, error) {
	entity := players.PlayerBioEntity{}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("atoi with %s returned error: %w", id, err)
	}
	playerInfo, err := a.client.PlayerInfo(intId)
	if err != nil {
		return players.PlayerBioEntity{}, fmt.Errorf("playerInfo with %v from %s returned error: %w", intId, reflect.TypeOf(a.client), err)
	}

	if len(playerInfo.Response) == 0 {
		return players.PlayerBioEntity{}, errors.New("empty player info response")
	}

	entity.FullName = playerInfo.Response[0].Name
	// api-basketball player response doesn't seem to have birth date in the example, only Age.
	// We might need to handle this.
	entity.BirthDate = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

	return entity, nil
}

func (a ApiBasketballStatsProviderAdapter) GetGamesStatsByPeriod(from, to time.Time) ([]games.GameStatEntity, error) {
	var passedGames []api_basketball.GameEntity

	// api-basketball typically filters by date (YYYY-MM-DD)
	current := from
	for !current.After(to) {
		dateStr := current.Format("2006-01-02")
		// Fetch ALL games for the date without league/season filters to comply with free plan
		response, err := a.client.Games(0, dateStr, "", "", "", "")
		if err != nil {
			return nil, fmt.Errorf("games for date %s returned error: %w", dateStr, err)
		}

		// Filter by leagueId in code
		for _, g := range response.Response {
			if g.League.Id == a.leagueId {
				passedGames = append(passedGames, g)
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	gamesStatsEntities := make([]games.GameStatEntity, 0, len(passedGames))

	for _, passedGame := range passedGames {
		// short status "FT" for finished games
		if passedGame.Status.Short != "FT" {
			continue
		}
		gameStatEntity := a.entityTransformer.TransformWithoutPlayers(passedGame)
		gameStatEntity.ExternalGameId = strconv.Itoa(passedGame.Id)
		gameStatEntity.HomeTeamExternalId = passedGame.Teams.Home.Id
		gameStatEntity.AwayTeamExternalId = passedGame.Teams.Away.Id

		gamesStatsEntities = append(gamesStatsEntities, gameStatEntity)
	}

	return gamesStatsEntities, nil
}

func (a ApiBasketballStatsProviderAdapter) EnrichGameStats(game games.GameStatEntity) (games.GameStatEntity, error) {
	if game.ExternalGameId == "" {
		return game, nil
	}

	gameId, err := strconv.Atoi(game.ExternalGameId)
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("atoi with %v returned error: %w", game.ExternalGameId, err)
	}

	playerStatsResponse, err := a.client.PlayersStatistics(gameId, 0, 0)
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("PlayersStatistics with %v from %s returned error: %w", gameId, reflect.TypeOf(a.client), err)
	}

	err = a.entityTransformer.MapPlayerStatistics(playerStatsResponse, game.HomeTeamExternalId, game.AwayTeamExternalId, &game)
	if err != nil {
		return games.GameStatEntity{}, fmt.Errorf("MapPlayerStatistics returned error: %w", err)
	}

	return game, nil
}

func NewApiBasketballStatsProviderAdapter(client api_basketball.ClientInterface, leagueId int) ApiBasketballStatsProviderAdapter {
	return ApiBasketballStatsProviderAdapter{
		client:            client,
		entityTransformer: api_basketball.NewEntityTransformer(),
		leagueId:          leagueId,
	}
}
