package handlers

import (
	"IMP/app/internal/abstract/custom_request"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"
)

// BindAndValidateRequestHandler is a helper function that binds and validates a request
// request must implement custom_request.CustomRequest interface
// controller handler could not have request parameter
func BindAndValidateRequestHandler[T custom_request.CustomRequest](controllerHandler func(w http.ResponseWriter, r T)) http.HandlerFunc {
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
		requestInstance.SetStorage(custom_request.NewCustomRequestStorage(
			mux.Vars(r),
			r.URL.Query().Encode(),
			decodedBody,
		))

		err := custom_request.Validate(requestInstance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		controllerHandler(w, requestInstance)
		//return bindAndValidate[custom_request.CustomRequest](controllerHandler)
	}
}
