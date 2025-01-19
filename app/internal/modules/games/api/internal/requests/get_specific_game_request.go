package requests

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type GetSpecificGameRequest struct {
	Id int
}

// Validate validates the {id} query-parameter to be integer
func (g *GetSpecificGameRequest) Validate(r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	g.Id = id
	return nil
}
