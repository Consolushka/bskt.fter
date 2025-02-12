package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/base/components/request_components"
)

// GetSpecificGameMetricsRequest expects request with:
//
// Path-parameters:
//   - id - game id
//
// Query-parameters:
//   - format - could be only 'pdf' or 'json'. Default is 'json'
//   - pers - comma-separated list of enums.ImpPERs
type GetSpecificGameMetricsRequest struct {
	custom_request.BaseRequest

	request_components.HasIdPathParam
	request_components.HasDateQueryParam
	request_components.HasPersQueryParam
}

func (g *GetSpecificGameMetricsRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	var parentValidators []func(storage custom_request.CustomRequestStorage) error

	for _, validator := range g.HasIdPathParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range g.HasDateQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range g.HasPersQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	return parentValidators
}
