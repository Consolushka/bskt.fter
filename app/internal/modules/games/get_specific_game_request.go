package games

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type getSpecificGameRequest struct {
	Id int
}

func (g *getSpecificGameRequest) Validate(r *http.Request) error {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}

	g.Id = id
	return nil
}
