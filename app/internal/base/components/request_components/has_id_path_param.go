package request_components

import (
	"IMP/app/internal/abstract"
	"strconv"
)

// HasIdPathParam expects integer id from path
type HasIdPathParam struct {
	abstract.BaseRequest

	id int
}

func (h *HasIdPathParam) Validators() []func(storage abstract.CustomRequestStorage) error {
	return []func(storage abstract.CustomRequestStorage) error{
		h.parseId,
	}
}

func (h *HasIdPathParam) Id() int {
	return h.id
}

func (h *HasIdPathParam) parseId(storage abstract.CustomRequestStorage) error {
	id, err := strconv.Atoi(storage.GetPathParam("id"))
	if err != nil {
		return err
	}
	h.id = id
	return nil
}
