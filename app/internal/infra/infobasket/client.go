package infobasket

import (
	"IMP/app/pkg/http"
	"fmt"
)

const (
	boxScoreEndpointPattern  = "/Widget/GetOnline/%v?format=json&lang=ru"
	teamGamesEndpointPattern = "/Widget/TeamGames/%v?format=json"
	scheduleEndpointPattern  = "/Comp/GetCalendar/?comps=%v&format=json"
)

type ClientInterface interface {
	BoxScore(gameId string) (GameBoxScoreResponse, error)
	ScheduledGames(compId int) ([]GameScheduleDto, error)
}

type Client struct {
	baseUrl string
}

func (c *Client) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	return http.Get[GameBoxScoreResponse](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)
}

func (c *Client) ScheduledGames(compId int) ([]GameScheduleDto, error) {
	return http.Get[[]GameScheduleDto](c.baseUrl+fmt.Sprintf(scheduleEndpointPattern, compId), nil)
}

func NewInfobasketClient() ClientInterface {
	return &Client{
		baseUrl: "https://reg.infobasket.su",
	}
}
