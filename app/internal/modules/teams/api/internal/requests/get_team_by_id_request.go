package requests

import (
	"IMP/app/internal/base/components/request_components"
)

// GetTeamByIdRequest expects request with:
//
// Path-parameters:
//   - id - game id: integer
type GetTeamByIdRequest struct {
	request_components.HasIdPathParam
}
