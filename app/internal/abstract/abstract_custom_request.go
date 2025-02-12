package abstract

type CustomRequest interface {
	Validators() []func(storage CustomRequestStorage) error
	GetStorage() *CustomRequestStorage
	SetStorage(storage CustomRequestStorage)
}
