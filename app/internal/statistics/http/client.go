package http

import (
	"encoding/json"
	"io"
	"net/http"
)

func Get[T any](url string, headers *map[string]string) T {
	var result T

	req, _ := http.NewRequest("GET", url, nil)

	if headers != nil {
		for key, value := range *headers {
			req.Header.Add(key, value)
		}
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
