package infobasket

import (
	"IMP/app/internal/abstract/http"
	"fmt"
)

const (
	boxScoreEndpointPattern  = "/GetOnline/%v?format=json&lang=ru"
	teamGamesEndpointPattern = "/TeamGames/%v?format=json"
)

type Client struct {
	baseUrl string
}

func (c *Client) BoxScore(gameId string) GameBoxScoreDto {
	result := http.Get[GameBoxScoreDto](c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result
}

func (c *Client) TeamGames(teamId string) TeamScheduleResponse {
	result := http.Get[TeamScheduleResponse](c.baseUrl+fmt.Sprintf(teamGamesEndpointPattern, teamId), nil)

	return result
}

func NewInfobasketClient() *Client {
	return &Client{
		baseUrl: "https://reg.infobasket.su/Widget",
	}
}
