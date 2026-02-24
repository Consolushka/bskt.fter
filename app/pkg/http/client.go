package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func Get[T any](url string, headers *map[string]string) (T, error) {
	return GetWithContext[T](context.Background(), url, headers)
}

func GetWithContext[T any](ctx context.Context, url string, headers *map[string]string) (T, error) {
	var result T

	if ctx == nil {
		return result, errors.New("nil context")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return result, fmt.Errorf("http.NewRequest with %s, %s, nil returned error: %w", "GET", url, err)
	}

	if headers != nil {
		for key, value := range *headers {
			req.Header.Add(key, value)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("http.DefaultClient.Do with %v returned error: %w", req, err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return result, fmt.Errorf("io.ReadAll with %s returned error: %w", url, err)
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return result, fmt.Errorf("http request to %s returned status %d: %s", url, res.StatusCode, string(body))
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, fmt.Errorf("json.Unmarshal with %v, %v returned error: %w", string(body), reflect.TypeOf(result), err)
	}

	return result, nil
}
