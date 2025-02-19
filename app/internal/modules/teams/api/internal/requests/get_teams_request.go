package requests

import (
	"IMP/app/internal/abstract/custom_request"
)

type GetTeamsRequest struct {
	custom_request.BaseRequest
}

func (g GetTeamsRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{}
}
