package middleware

import (
	"IMP/app/internal/utils/request_utils"
	"net/http"
)

func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := request_utils.ParseQueryParams(r)
		var contentType string

		switch queryParams["format"] {
		case "pdf":
			contentType = "application/pdf"
		default:
			contentType = "application/json"
		}

		w.Header().Set("Content-Type", contentType)

		next.ServeHTTP(w, r)
	})
}
