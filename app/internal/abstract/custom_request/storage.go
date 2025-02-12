package custom_request

import (
	"strings"
)

type CustomRequestStorage struct {
	pathParams  map[string]string
	queryParams map[string]string
	body        map[string]interface{}
}

func NewCustomRequestStorage(pathParams map[string]string, encodedQueryParams string, body map[string]interface{}) CustomRequestStorage {
	var queryParams map[string]string
	if encodedQueryParams != "" {
		queryParams = make(map[string]string)
		queryParamsPairs := strings.Split(encodedQueryParams, "&")
		for _, pair := range queryParamsPairs {
			parts := strings.Split(pair, "=")

			queryParams[parts[0]] = parts[1]
		}
	}

	return CustomRequestStorage{
		pathParams:  pathParams,
		queryParams: queryParams,
		body:        body,
	}
}

func (c *CustomRequestStorage) GetPathParam(key string) string {
	return c.pathParams[key]
}

func (c *CustomRequestStorage) GetQueryParam(key string) string {
	return c.queryParams[key]
}

func (c *CustomRequestStorage) GetBodyParam(key string) (interface{}, bool) {
	value, exists := c.body[key]
	return value, exists
}
