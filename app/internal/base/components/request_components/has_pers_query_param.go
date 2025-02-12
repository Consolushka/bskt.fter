package request_components

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/modules/imp/domain/enums"
	"strings"
)

// HasPersQueryParam expects pers from enums.ImpPERs split by ','
type HasPersQueryParam struct {
	abstract.BaseRequest

	pers []enums.ImpPERs
}

func (h *HasPersQueryParam) Validators() []func(storage abstract.CustomRequestStorage) error {
	return []func(storage abstract.CustomRequestStorage) error{
		h.parsePERS,
	}
}

func (h *HasPersQueryParam) parsePERS(storage abstract.CustomRequestStorage) error {
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
