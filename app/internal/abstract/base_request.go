package abstract

type BaseRequest struct {
	storage CustomRequestStorage
}

func (b *BaseRequest) SetStorage(storage CustomRequestStorage) {
	b.storage = storage
}

func (b *BaseRequest) GetStorage() *CustomRequestStorage {
	return &b.storage
}
