package requests

import (
	"IMP/app/internal/abstract"
	"strconv"
)

// GetTeamByIdRequest expects request with:
//
// Path-parameters:
//   - id - game id: integer
type GetTeamByIdRequest struct {
	abstract.BaseRequest

	id int
}

// Validate validates the {id} path-parameter to be integer
func (g *GetTeamByIdRequest) Validate() error {
	return g.parseAll(
		g.parseId,
	)
}

func (g *GetTeamByIdRequest) Id() int {
	return g.id
}

func (g *GetTeamByIdRequest) parseAll(parsers ...func() error) error {
	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetTeamByIdRequest) parseId() error {
	id, err := strconv.Atoi(g.GetStorage().GetPathParam("id"))
	if err != nil {
		return err
	}
	g.id = id
	return nil
}
