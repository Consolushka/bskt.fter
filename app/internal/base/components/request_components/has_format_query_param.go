package request_components

import (
	"IMP/app/internal/abstract"
	"errors"
)

// HasFormatQueryParam expects format from query 'json' or 'pdf'. Default is 'json'
type HasFormatQueryParam struct {
	abstract.BaseRequest

	format string
}

func (h *HasFormatQueryParam) Validators() []func(storage abstract.CustomRequestStorage) error {
	return []func(storage abstract.CustomRequestStorage) error{
		h.parseFormat,
	}
}

func (h *HasFormatQueryParam) Format() string {
	return h.format
}

func (h *HasFormatQueryParam) parseFormat(storage abstract.CustomRequestStorage) error {
	format := storage.GetQueryParam("format")
	if format == "" {
		h.format = "json"
		return nil
	}

	if format != "json" && format != "pdf" {
		return errors.New("invalid format")
	}
	h.format = format
	return nil
}
