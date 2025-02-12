package request_components

import (
	"IMP/app/internal/abstract/custom_request"
	"strconv"
)

// HasIdPathParam expects integer id from path
type HasIdPathParam struct {
	custom_request.BaseRequest

	id int
}

func (h *HasIdPathParam) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{
		h.parseId,
	}
}

func (h *HasIdPathParam) Id() int {
	return h.id
}

func (h *HasIdPathParam) parseId(storage custom_request.CustomRequestStorage) error {
	id, err := strconv.Atoi(storage.GetPathParam("id"))
	if err != nil {
		return err
	}
	h.id = id
	return nil
}
