package request_components

import (
	"IMP/app/internal/abstract/custom_request"
	"errors"
	"time"
)

// HasDateQueryParam expects date in format 'dd-mm-yyyy'
type HasDateQueryParam struct {
	custom_request.BaseRequest

	date *time.Time
}

func (h *HasDateQueryParam) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{
		h.parseDate,
	}
}

func (h *HasDateQueryParam) Date() *time.Time {
	return h.date
}

func (h *HasDateQueryParam) parseDate(storage custom_request.CustomRequestStorage) error {
	queryDate := storage.GetQueryParam("date")

	if queryDate == "" {
		return errors.New("query-parameter 'date' is required")
	}

	date, err := time.Parse("02-01-2006", queryDate)
	if err != nil {
		return errors.New("'date' parameter is not a valid date. valid format: dd-mm-yyyy")
	}

	h.date = &date
	return nil
}
