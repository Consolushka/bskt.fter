package requests

import (
	"IMP/app/internal/abstract"
)

type GetTeamsRequest struct {
	abstract.BaseRequest
}

func (g *GetTeamsRequest) Validate() error {
	return g.parseAll()
}

func (g *GetTeamsRequest) parseAll(parsers ...func() error) error {
	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}
	return nil
}
