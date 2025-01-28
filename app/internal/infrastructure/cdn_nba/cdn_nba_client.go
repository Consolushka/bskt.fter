package cdn_nba

import (
	"IMP/app/internal/abstract"
	"fmt"
)

const boxScoreEndpointPattern = "/liveData/boxscore/boxscore_%v.json"
const fullSeasonEndpoint = "/staticData/scheduleLeagueV2_14.json"

type Client struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func NewCdnNbaClient() *Client {
	return &Client{
		baseUrl:    "https://cdn.nba.com/static/json",
		httpClient: abstract.NewHttpClient(),
	}
}

func (c Client) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result["game"].(map[string]interface{})
}

func (c Client) ScheduleSeason() map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fullSeasonEndpoint, nil)

	return result["leagueSchedule"].(map[string]interface{})
}
