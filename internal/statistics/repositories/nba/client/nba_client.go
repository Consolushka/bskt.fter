package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const boxScoreEndpointPattern = "/boxscore/boxscore_%v.json"

type NbaClient struct {
	baseUrl string
}

func NewNbaClient() *NbaClient {
	return &NbaClient{
		baseUrl: "https://cdn.nba.com/static/json/liveData",
	}
}

func (c NbaClient) BoxScore(gameId string) map[string]interface{} {
	var result map[string]interface{}

	url := c.baseUrl + fmt.Sprintf(boxScoreEndpointPattern, gameId)

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

	return result["game"].(map[string]interface{})
}
