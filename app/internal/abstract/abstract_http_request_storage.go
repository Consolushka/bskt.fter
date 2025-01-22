package abstract

type CustomRequestStorage struct {
	queryParams map[string]string
	body        map[string]interface{}
}

func NewCustomRequestStorage(queryParams map[string]string, body map[string]interface{}) CustomRequestStorage {
	return CustomRequestStorage{
		queryParams: queryParams,
		body:        body,
	}
}

func (c *CustomRequestStorage) GetQueryParam(key string) string {
	return c.queryParams[key]
}

func (c *CustomRequestStorage) GetBodyParam(key string) (interface{}, bool) {
	value, exists := c.body[key]
	return value, exists
}
