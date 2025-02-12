package requests

import "IMP/app/internal/base/components/request_components"

// GetTeamGamesRequest expects request with:
//
// Path-parameters:
//   - id - game id: integer
type GetTeamGamesRequest struct {
	request_components.HasIdPathParam
}
