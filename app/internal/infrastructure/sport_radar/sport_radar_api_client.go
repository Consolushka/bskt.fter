package sport_radar

import (
	"IMP/app/internal/abstract"
)

const (
	gameSummaryEndpoint = "summary.json"
)

type Client struct {
	baseUrl string
	version string
	lang    string
	apiKey  string

	httpClient *abstract.HttpClient
}

func NewSportRadarApiClient() *Client {
	return &Client{
		baseUrl:    "https://api.sportradar.com/nba/trial/v8",
		version:    "v8",
		lang:       "en",
		apiKey:     "piUTvn6SPhj5EX8NIS9vOxHDGKRaMNwYLXVD5u9O",
		httpClient: abstract.NewHttpClient(),
	}
}

func (c Client) GameSummary(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+"/"+c.lang+"/games/"+gameId+"/"+gameSummaryEndpoint+"?api_key="+c.apiKey, nil)

	return result.(map[string]interface{})
}
