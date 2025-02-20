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
}

func NewSportRadarApiClient() *Client {
	return &Client{
		baseUrl: "https://api.sportradar.com/nba/trial/v8",
		version: "v8",
		lang:    "en",
		apiKey:  os.Getenv("SPORTRADAR_API_KEY"),
	}
}

func (c Client) GameSummary(gameId string) GameBoxScoreDTO {
	result := http.Get[GameBoxScoreDTO](c.baseUrl+"/"+c.lang+"/games/"+gameId+"/"+gameSummaryEndpoint+"?api_key="+c.apiKey, nil)

	return result
}
