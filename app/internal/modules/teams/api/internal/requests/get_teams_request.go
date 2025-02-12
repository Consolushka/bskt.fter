package requests

import (
	"IMP/app/internal/abstract"
)

type GetTeamsRequest struct {
	abstract.BaseRequest
}

func (g GetTeamsRequest) Validators() []func(storage abstract.CustomRequestStorage) error {
	return []func(storage abstract.CustomRequestStorage) error{}
}
