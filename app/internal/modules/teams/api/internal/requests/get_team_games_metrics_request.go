package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/base/components/request_components"
)

// GetTeamByIdGamesMetricsRequest expects request with:
//
// Path-parameters:
//   - id - team id
//
// Query-parameters:
//   - pers - comma-separated list of PERS
type GetTeamByIdGamesMetricsRequest struct {
	custom_request.BaseRequest

	request_components.HasIdPathParam
	request_components.HasPersQueryParam
}

func (g *GetTeamByIdGamesMetricsRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	var parentValidators []func(storage custom_request.CustomRequestStorage) error

	for _, validator := range g.HasIdPathParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range g.HasPersQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	return parentValidators
}
