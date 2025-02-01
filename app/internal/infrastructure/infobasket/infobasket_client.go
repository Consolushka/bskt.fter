package infobasket

import (
	"IMP/app/internal/abstract"
	"fmt"
)

const (
	boxScoreEndpointPattern  = "/GetOnline/%v?format=json&lang=ru"
	teamGamesEndpointPattern = "/TeamGames/%v?format=json"
)

type Client struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func (c *Client) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result.(map[string]interface{})
}

func (c *Client) TeamGames(teamId string) []map[string]interface{} {
	rawResult := c.httpClient.Get(c.baseUrl+fmt.Sprintf(teamGamesEndpointPattern, teamId), nil)

	interfaceSlice := rawResult.([]interface{})
	result := make([]map[string]interface{}, len(interfaceSlice))

	for i, v := range interfaceSlice {
		result[i] = v.(map[string]interface{})
	}

	return result
}

func NewInfobasketClient() *Client {
	return &Client{
		baseUrl:    "https://reg.infobasket.su/Widget",
		httpClient: abstract.NewHttpClient(),
	}
}
