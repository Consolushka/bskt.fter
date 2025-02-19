package request_utils

import (
	"net/http"
	"strings"
)

func ParseQueryParams(r *http.Request) map[string]string {
	var queryParams map[string]string
	encodedQueryParams := r.URL.Query().Encode()
	if encodedQueryParams != "" {
		queryParams = make(map[string]string)
		queryParamsPairs := strings.Split(encodedQueryParams, "&")
		for _, pair := range queryParamsPairs {
			parts := strings.Split(pair, "=")

			queryParams[parts[0]] = parts[1]
		}

	}

	return queryParams
}
