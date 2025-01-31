package handlers

import (
	"IMP/app/internal/abstract"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

func BindAndValidateRequestHandler[T abstract.CustomRequest](controllerHandler func(w http.ResponseWriter, r T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var decodedBody map[string]interface{}

		var request T

		if r.Body != nil && r.ContentLength > 0 {
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(&decodedBody); err != nil {
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()
		} else {
			// Initialize empty map if body is nil or empty
			decodedBody = make(map[string]interface{})
		}

		requestInstance := reflect.New(reflect.TypeOf(request).Elem()).Interface().(T)
		requestInstance.SetStorage(abstract.NewCustomRequestStorage(
			mux.Vars(r),
			r.URL.Query().Encode(),
			decodedBody,
		))

		err := requestInstance.Validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		controllerHandler(w, requestInstance)
	}
}
