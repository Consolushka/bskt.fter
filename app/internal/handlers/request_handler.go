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

func bindAndValidate[T custom_request.CustomRequest](controllerHandler interface{}) http.HandlerFunc {
	handlerType := reflect.TypeOf(controllerHandler)

	// Check if it's a function
	if handlerType.Kind() != reflect.Func {
		panic("handler must be a function")
	}

	// Get number of parameters
	numParams := handlerType.NumIn()

	return func(w http.ResponseWriter, r *http.Request) {
		if numParams == 1 {
			// Simple handler with just ResponseWriter
			reflect.ValueOf(controllerHandler).Call([]reflect.Value{
				reflect.ValueOf(w),
			})
			return
		}

		// Handler with request binding
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

		reflect.ValueOf(controllerHandler).Call([]reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(requestInstance),
		})
	}
}
