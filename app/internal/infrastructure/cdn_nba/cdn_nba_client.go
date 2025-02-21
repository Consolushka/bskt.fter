package cdn_nba

import (
	"IMP/app/internal/abstract/http"
	"fmt"
)

const boxScoreEndpointPattern = "/liveData/boxscore/boxscore_%v.json"
const fullSeasonEndpoint = "/staticData/scheduleLeagueV2_14.json"

type Client struct {
	baseUrl string
}

func NewCdnNbaClient() *Client {
	return &Client{
		baseUrl: "https://cdn.nba.com/static/json",
	}
}

func (c Client) BoxScore(gameId string) BoxScoreDto {
	result := http.Get[GameBoxScoreResponse](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result.Game
}

// todo: save to static file
func (c Client) ScheduleSeason() SeasonScheduleDto {
	result := http.Get[ScheduleResponse](c.baseUrl+fullSeasonEndpoint, nil)

	return result.Schedule
}
