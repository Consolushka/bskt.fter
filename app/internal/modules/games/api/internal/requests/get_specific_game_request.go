package requests

import (
	"IMP/app/internal/base/components/request_components"
)

// GetSpecificGameRequest expects request with:
//
// Path-parameters:
//   - id - game id: integer
type GetSpecificGameRequest struct {
	request_components.HasIdPathParam
}
