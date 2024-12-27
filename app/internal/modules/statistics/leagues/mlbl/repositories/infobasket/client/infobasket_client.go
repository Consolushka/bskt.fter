package client

import (
	"IMP/app/internal/abstract"
	"fmt"
)

const (
	boxScoreEndpointPattern = "/GetOnline/%v?format=json&lang=ru"
)

type InfobasketClient struct {
	baseUrl string

	httpClient *abstract.HttpClient
}

func (c InfobasketClient) BoxScore(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl + fmt.Sprintf(boxScoreEndpointPattern, gameId))

	return result
}

func NewInfobasketClient() *InfobasketClient {
	return &InfobasketClient{
		baseUrl:    "https://reg.infobasket.su/Widget",
		httpClient: abstract.NewHttpClient(),
	}
}
