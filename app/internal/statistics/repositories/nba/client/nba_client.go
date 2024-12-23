package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const boxScoreEndpointPattern = "/boxscore/boxscore_%v.json"
const todaysGamesEndpoint = "/scoreboard/todaysScoreboard_00.json"

type NbaClient struct {
	baseUrl string
}

func NewNbaClient() *NbaClient {
	return &NbaClient{
		baseUrl: "https://cdn.nba.com/static/json/liveData",
	}
}

func (c NbaClient) BoxScore(gameId string) map[string]interface{} {
	result := get(c.baseUrl + fmt.Sprintf(boxScoreEndpointPattern, gameId))

	return result["game"].(map[string]interface{})
}

func (c NbaClient) TodaysGames() map[string]interface{} {
	result := get(c.baseUrl + todaysGamesEndpoint)

	return result["scoreboard"].(map[string]interface{})
}

func get(url string) map[string]interface{} {
	var result map[string]interface{}

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
