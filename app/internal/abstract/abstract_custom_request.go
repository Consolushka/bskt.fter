package abstract

import "net/http"

type CustomRequest interface {
	Validate(r *http.Request) error
	GetStorage() *CustomRequestStorage
	SetStorage(storage CustomRequestStorage)
}
