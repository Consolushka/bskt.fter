package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/base/components/request_components"
)

type PlayerGamesMetricsRequest struct {
	custom_request.BaseRequest

	request_components.HasIdPathParam
	request_components.HasPersQueryParam
}

func (h *PlayerGamesMetricsRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	var parentValidators []func(storage custom_request.CustomRequestStorage) error

	for _, validator := range h.HasIdPathParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range h.HasPersQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	return parentValidators
}
