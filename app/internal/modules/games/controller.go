package games

import (
	"encoding/json"
	"net/http"
)

type controller struct {
	service *Service
}

func newController() *controller {
	return &controller{
		service: NewService(),
	}
}

// getGame
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *controller) getGame(w http.ResponseWriter, r *getSpecificGameRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetGame(r.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fromGameModel(gameModel)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// getGameMetrics
//
// retrieve specific game and then calculate IMP metrics for every player
func (c *controller) getGameMetrics(w http.ResponseWriter, r *getSpecificGameRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetGameMetrics(r.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fromGameMetricsModel(gameModel)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
