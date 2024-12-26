package abstract

import (
	"encoding/json"
	"io"
	"net/http"
)

type HttpClient struct{}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}

// Get send GET request to given url w/o body and headers (yet)
func (c *HttpClient) Get(url string) map[string]interface{} {
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
