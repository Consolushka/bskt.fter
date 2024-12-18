package client

import (
	"encoding/json"
	"io"
	"net/http"
)

const (
	gameSummaryEndpoint = "summary.json"
)

type SportRadarClient struct {
	baseUrl string
	version string
	lang    string
	apiKey  string
}

func NewSportRadarClient() *SportRadarClient {
	return &SportRadarClient{
		baseUrl: "https://api.sportradar.com/nba/trial/v8",
		version: "v8",
		lang:    "en",
		apiKey:  "piUTvn6SPhj5EX8NIS9vOxHDGKRaMNwYLXVD5u9O",
	}
}

func (c SportRadarClient) GameSummary(gameId string) map[string]interface{} {
	var result map[string]interface{}

	url := c.baseUrl + "/" + c.lang + "/games/" + gameId + "/" + gameSummaryEndpoint + "?api_key=" + c.apiKey

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)

	err := json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	return result
}
