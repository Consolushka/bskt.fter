package requests

import "IMP/app/internal/abstract/custom_request"

type GetLeaguesRequest struct {
	custom_request.BaseRequest
}

func (g GetLeaguesRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	return []func(storage custom_request.CustomRequestStorage) error{}
}
