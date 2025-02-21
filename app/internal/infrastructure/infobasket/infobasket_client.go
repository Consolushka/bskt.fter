package infobasket

import (
	"IMP/app/internal/abstract/http"
	"fmt"
)

const (
	boxScoreEndpointPattern  = "/Widget/GetOnline/%v?format=json&lang=ru"
	teamGamesEndpointPattern = "/Widget/TeamGames/%v?format=json"
	scheduleEndpointPattern  = "/Comp/GetCalendar/?comps=%v&format=json"
)

type Client struct {
	baseUrl string
}

func (c *Client) BoxScore(gameId string) GameBoxScoreResponse {
	result := http.Get[GameBoxScoreResponse](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result
}

func (c *Client) TeamGames(teamId string) TeamScheduleResponse {
	result := http.Get[TeamScheduleResponse](c.baseUrl+fmt.Sprintf(teamGamesEndpointPattern, teamId), nil)

	return result
}

func (c *Client) ScheduledGames(compId int) []GameScheduleDto {
	result := http.Get[[]GameScheduleDto](c.baseUrl+fmt.Sprintf(scheduleEndpointPattern, compId), nil)

	return result
}

func NewInfobasketClient() *Client {
	return &Client{
		baseUrl: "https://reg.infobasket.su",
	}
}
