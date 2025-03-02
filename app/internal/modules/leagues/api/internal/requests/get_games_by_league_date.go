package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/base/components/request_components"
)

type GetGamesByLeagueAndDate struct {
	custom_request.BaseRequest

	request_components.HasIdPathParam
	request_components.HasDateQueryParam
}

func (g *GetGamesByLeagueAndDate) Validators() []func(storage custom_request.CustomRequestStorage) error {
	var parentValidators []func(storage custom_request.CustomRequestStorage) error

	for _, validator := range g.HasIdPathParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range g.HasDateQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	return parentValidators
}
