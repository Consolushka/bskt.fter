package games

import (
	"encoding/json"
	"net/http"
)

type controller struct {
	service *service
}

func newController() *controller {
	return &controller{
		service: newService(),
	}
}

// getGame handles GET /game/{id}
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *controller) getGame(w http.ResponseWriter, r *getSpecificGameRequest) {
	w.Header().Set("Content-Type", "application/json")

	stats, err := c.service.getGame(r.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(stats)
}
