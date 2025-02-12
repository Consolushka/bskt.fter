package requests

import "IMP/app/internal/base/components/request_components"

// GetGamesRequest expects request with:
//
// Query-parameters:
//   - date - date in format dd-mm-yyyy
type GetGamesRequest struct {
	request_components.HasDateQueryParam
}
