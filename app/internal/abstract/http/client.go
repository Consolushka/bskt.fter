package http

import (
	"encoding/json"
	"io"
	"net/http"
)

type Client struct{}

func NewHttpClient() *Client {
	return &Client{}
}

func Get[T any](url string, apiKey *string) T {
	var result T

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	if apiKey != nil {
		req.Header.Add("Authorization", *apiKey)
	}

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

// Get send GET request to given url w/o body and headers (yet)
//
//todo:use better abstraction
func (c *Client) Get(url string, apiKey *string) interface{} {
	var result interface{}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	if apiKey != nil {
		req.Header.Add("Authorization", *apiKey)
	}

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
