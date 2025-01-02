package infobasket

import (
	"IMP/app/internal/abstract"
	"fmt"
)

const (
	boxScoreEndpointPattern = "/GetOnline/%v?format=json&lang=ru"
)

type Client struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func (c Client) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+fmt.Sprintf(boxScoreEndpointPattern, gameId), nil)

	return result
}

func NewInfobasketClient() *Client {
	return &Client{
		baseUrl:    "https://reg.infobasket.su/Widget",
		httpClient: abstract.NewHttpClient(),
	}
}
