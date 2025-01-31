package requests

import (
	"IMP/app/internal/abstract"
	"errors"
	"time"
)

// GetGamesRequest expects request with:
//
// Query-parameters:
//   - date - date in format dd-mm-yyyy
type GetGamesRequest struct {
	abstract.BaseRequest

	date *time.Time
}

// Validate validates the {id} path-parameter to be integer
func (g *GetGamesRequest) Validate() error {
	return g.parseAll(
		g.parseDate,
	)
}

func (g *GetGamesRequest) Date() *time.Time {
	return g.date
}

func (g *GetGamesRequest) parseAll(parsers ...func() error) error {
	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetGamesRequest) parseDate() error {
	queryDate := g.GetStorage().GetQueryParam("date")

	if queryDate == "" {
		return errors.New("query-parameter 'date' is required")
	}

	date, err := time.Parse("02-01-2006", queryDate)
	if err != nil {
		return errors.New("'date' parameter is not a valid date. valid format: dd-mm-yyyy")
	}

	g.date = &date
	return nil
}
