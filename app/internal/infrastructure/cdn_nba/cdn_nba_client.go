package cdn_nba

import (
	"IMP/app/internal/abstract/http"
	"fmt"
)

const boxScoreEndpointPattern = "/liveData/boxscore/boxscore_%v.json"
const fullSeasonEndpoint = "/staticData/scheduleLeagueV2_14.json"

type Client struct {
	baseUrl string

	httpClient *http.Client
}

func NewCdnNbaClient() *Client {
	return &Client{
		baseUrl:    "https://cdn.nba.com/static/json",
		httpClient: http.NewHttpClient(),
	}
}

func (c Client) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result.(map[string]interface{})["game"].(map[string]interface{})
}

func (c Client) ScheduleSeason() map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fullSeasonEndpoint, nil)

	return result.(map[string]interface{})["leagueSchedule"].(map[string]interface{})
}
