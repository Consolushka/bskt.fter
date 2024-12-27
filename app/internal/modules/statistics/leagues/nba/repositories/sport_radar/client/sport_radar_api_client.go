package client

import (
	"IMP/app/internal/abstract"
)

const (
	gameSummaryEndpoint = "summary.json"
)

type SportRadarApiClient struct {
	baseUrl string
	version string
	lang    string
	apiKey  string

	httpClient *abstract.HttpClient
}

func NewSportRadarApiClient() *SportRadarApiClient {
	return &SportRadarApiClient{
		baseUrl:    "https://api.sportradar.com/nba/trial/v8",
		version:    "v8",
		lang:       "en",
		apiKey:     "piUTvn6SPhj5EX8NIS9vOxHDGKRaMNwYLXVD5u9O",
		httpClient: abstract.NewHttpClient(),
	}
}

func (c SportRadarApiClient) GameSummary(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl + "/" + c.lang + "/games/" + gameId + "/" + gameSummaryEndpoint + "?api_key=" + c.apiKey)

	return result
}
