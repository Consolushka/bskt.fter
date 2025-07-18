package cdn_nba

import (
	"IMP/app/pkg/http"
	"fmt"
)

const boxScoreEndpointPattern = "/liveData/boxscore/boxscore_%v.json"
const fullSeasonEndpoint = "/staticData/scheduleLeagueV2_14.json"

type ClientInterface interface {
	BoxScore(gameId string) (GameBoxScoreResponse, error)
	ScheduleSeason() (ScheduleResponse, error)
}

type Client struct {
	baseUrl string
}

func NewCdnNbaClient() *Client {
	return &Client{
		baseUrl: "https://cdn.nba.com/static/json",
	}
}

func (c Client) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	return http.Get[GameBoxScoreResponse](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)
}

func (c Client) ScheduleSeason() (ScheduleResponse, error) {
	return http.Get[ScheduleResponse](c.baseUrl+fullSeasonEndpoint, nil)
}
