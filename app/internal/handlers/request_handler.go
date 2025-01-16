package handlers

import (
	"IMP/app/internal/abstract"
	"net/http"
	"reflect"
)

func BindAndValidateRequestHandler[T abstract.CustomRequest](controllerHandler func(w http.ResponseWriter, r T)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request T

		requestInstance := reflect.New(reflect.TypeOf(request).Elem()).Interface().(T)

		err := requestInstance.Validate(r)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		controllerHandler(w, requestInstance)
	}
}
