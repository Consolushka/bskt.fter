package client

import (
	"FTER/app/internal/abstract"
	"fmt"
)

const boxScoreEndpointPattern = "/boxscore/boxscore_%v.json"
const todaysGamesEndpoint = "/scoreboard/todaysScoreboard_00.json"

type NbaComApiClient struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func NewNbaComApiClient() *NbaComApiClient {
	return &NbaComApiClient{
		baseUrl:    "https://cdn.nba.com/static/json/liveData",
		httpClient: abstract.NewHttpClient(),
	}
}

func (c NbaComApiClient) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl + fmt.Sprintf(boxScoreEndpointPattern, gameId))

	return result["game"].(map[string]interface{})
}

func (c NbaComApiClient) TodaysGames() map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl + todaysGamesEndpoint)

	return result["scoreboard"].(map[string]interface{})
}
