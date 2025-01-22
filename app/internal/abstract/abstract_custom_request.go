package abstract

type CustomRequest interface {
	Validate() error
	GetStorage() *CustomRequestStorage
	SetStorage(storage CustomRequestStorage)
}
