package cdn_nba

import (
	"IMP/app/internal/abstract"
	"fmt"
)

const boxScoreEndpointPattern = "/boxscore/boxscore_%v.json"
const todaysGamesEndpoint = "/scoreboard/todaysScoreboard_00.json"

type Client struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func NewCdnNbaClient() *Client {
	return &Client{
		baseUrl:    "https://cdn.nba.com/static/json/liveData",
		httpClient: abstract.NewHttpClient(),
	}
}

func (c Client) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result["game"].(map[string]interface{})
}

func (c Client) TodaysGames() map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+todaysGamesEndpoint, nil)

	return result["scoreboard"].(map[string]interface{})
}
