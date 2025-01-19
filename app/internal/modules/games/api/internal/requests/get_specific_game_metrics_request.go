package requests

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type GetSpecificGameMetricsRequest struct {
	id     int
	format string
}

// Validate validates the {id} query-parameter and {format} parameter (could be only 'pdf' or 'json')
func (g *GetSpecificGameMetricsRequest) Validate(r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	format := vars["format"]
	if format != "json" && format != "pdf" {
		return errors.New("invalid format")
	}

	g.id = id
	g.format = format
	return nil
}

func (g *GetSpecificGameMetricsRequest) Id() int {
	return g.id
}

func (g *GetSpecificGameMetricsRequest) Format() string {
	return g.format
}
