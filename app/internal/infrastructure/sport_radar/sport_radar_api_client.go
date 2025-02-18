package sport_radar

import (
	"IMP/app/internal/abstract/http"
	"os"
)

const (
	gameSummaryEndpoint = "summary.json"
)

type Client struct {
	baseUrl string
	version string
	lang    string
	apiKey  string

	httpClient *http.Client
}

func NewSportRadarApiClient() *Client {
	return &Client{
		baseUrl:    "https://api.sportradar.com/nba/trial/v8",
		version:    "v8",
		lang:       "en",
		apiKey:     os.Getenv("SPORTRADAR_API_KEY"),
		httpClient: http.NewHttpClient(),
	}
}

func (c Client) GameSummary(gameId string) map[string]interface{} {
	result := c.httpClient.Get(c.baseUrl+"/"+c.lang+"/games/"+gameId+"/"+gameSummaryEndpoint+"?api_key="+c.apiKey, nil)

	return result.(map[string]interface{})
}
