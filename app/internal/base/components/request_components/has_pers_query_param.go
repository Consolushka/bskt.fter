package request_components

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/modules/imp/domain/enums"
	"strings"
)

// HasPersQueryParam expects pers from enums.ImpPERs split by ','
type HasPersQueryParam struct {
	custom_request.BaseRequest

	pers []enums.ImpPERs
}

func (h *HasPersQueryParam) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{
		h.parsePERS,
	}
}

func (h *HasPersQueryParam) parsePERS(storage custom_request.CustomRequestStorage) error {
	persInterface := storage.GetQueryParam("pers")
	if persInterface != "" {
		persArray := strings.Split(persInterface, "%2C")
		h.pers = make([]enums.ImpPERs, len(persArray))
		for i, v := range persArray {
			h.pers[i] = enums.ImpPERs(v)
		}
	}
	return nil
}

func (h *HasPersQueryParam) Pers() []enums.ImpPERs {
	return h.pers
}
