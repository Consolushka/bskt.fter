package requests

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/modules/imp/domain/enums"
	"errors"
	"strconv"
	"strings"
)

// GetSpecificGameMetricsRequest expects request with:
//
// Path-parameters:
//   - id - game id
//
// Query-parameters:
//   - format - could be only 'pdf' or 'json'. Default is 'json'
//   - pers - comma-separated list of PERS
type GetSpecificGameMetricsRequest struct {
	abstract.BaseRequest

	id int

	format string
	pers   []enums.ImpPERs
}

// Validate validates the {id} path-parameter and {format} parameter (could be only 'pdf' or 'json')
func (g *GetSpecificGameMetricsRequest) Validate() error {
	return g.parseAll(
		g.parseId,
		g.parseFormat,
		g.parsePERS,
	)
}

func (g *GetSpecificGameMetricsRequest) parseAll(parsers ...func() error) error {
	for _, parser := range parsers {
		if err := parser(); err != nil {
			return err
		}
	}
	return nil
}

func (g *GetSpecificGameMetricsRequest) parseId() error {
	id, err := strconv.Atoi(g.GetStorage().GetPathParam("id"))
	if err != nil {
		return err
	}
	g.id = id
	return nil
}

func (g *GetSpecificGameMetricsRequest) parseFormat() error {
	format := g.GetStorage().GetQueryParam("format")
	if format == "" {
		g.format = "json"
		return nil
	}

	if format != "json" && format != "pdf" {
		return errors.New("invalid format")
	}
	g.format = format
	return nil
}

func (g *GetSpecificGameMetricsRequest) parsePERS() error {
	persInterface := g.GetStorage().GetQueryParam("pers")
	if persInterface != "" {
		persArray := strings.Split(persInterface, "%2C")
		g.pers = make([]enums.ImpPERs, len(persArray))
		for i, v := range persArray {
			g.pers[i] = enums.ImpPERs(v)
		}
	}
	return nil
}

func (g *GetSpecificGameMetricsRequest) Id() int {
	return g.id
}

func (g *GetSpecificGameMetricsRequest) Format() string {
	return g.format
}

func (g *GetSpecificGameMetricsRequest) Pers() []enums.ImpPERs {
	return g.pers
}
