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

// Validate validates for properties of inherited request
func Validate(request CustomRequest) error {
	return parseAll(
		request.Validators(),
		*request.GetStorage(),
	)
}

func parseAll(parsers []func(storage CustomRequestStorage) error, storage CustomRequestStorage) error {
	for _, parser := range parsers {
		if err := parser(storage); err != nil {
			return err
		}
	}
	return nil
}
