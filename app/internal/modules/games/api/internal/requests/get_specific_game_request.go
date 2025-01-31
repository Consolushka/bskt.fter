package requests

import (
	"IMP/app/internal/abstract"
	"strconv"
)

// GetSpecificGameRequest expects request with:
//
// Path-parameters:
//   - id - game id: integer
type GetSpecificGameRequest struct {
	abstract.BaseRequest

	id int
}

// Validate validates the {id} path-parameter to be integer
func (g *GetSpecificGameRequest) Validate() error {
	return g.parseAll(
		g.parseId,
	)
}

func (g *GetSpecificGameRequest) Id() int {
	return g.id
}

func (g *GetSpecificGameRequest) parseAll(parsers ...func() error) error {
	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetSpecificGameRequest) parseId() error {
	id, err := strconv.Atoi(g.GetStorage().GetPathParam("id"))
	if err != nil {
		return err
	}
	g.id = id
	return nil
}
