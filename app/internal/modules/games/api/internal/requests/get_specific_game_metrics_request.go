package requests

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/modules/imp/domain/enums"
	"errors"
	"strconv"
)

type GetSpecificGameMetricsRequest struct {
	abstract.BaseRequest

	id int

	format string
	pers   []enums.ImpPERs
}

// Validate validates the {id} query-parameter and {format} parameter (could be only 'pdf' or 'json')
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
	id, err := strconv.Atoi(g.GetStorage().GetQueryParam("id"))
	if err != nil {
		return err
	}
	g.id = id
	return nil
}

func (g *GetSpecificGameMetricsRequest) parseFormat() error {
	format, exists := g.GetStorage().GetBodyParam("format")
	if !exists {
		g.format = "json"
		return nil
	}
	if format != "json" && format != "pdf" {
		return errors.New("invalid format")
	}
	g.format = format.(string)
	return nil
}

func (g *GetSpecificGameMetricsRequest) parsePERS() error {
	persInterface, exists := g.GetStorage().GetBodyParam("pers")
	if exists {
		persArray := persInterface.([]interface{})
		g.pers = make([]enums.ImpPERs, len(persArray))
		for i, v := range persArray {
			g.pers[i] = enums.ImpPERs(v.(string))
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
