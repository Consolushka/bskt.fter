package internal

import (
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/games/api/internal/formatters"
	"IMP/app/internal/modules/games/api/internal/requests"
	"IMP/app/internal/modules/games/resources"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller struct {
	service *games.Service

	logger *logrus.Logger
}

func NewController() *Controller {
	return &Controller{
		service: games.NewService(),
	}
}

// GetGames returns all games filtered by date
func (c *Controller) GetGames(w http.ResponseWriter, r *requests.GetGamesRequest) {
	var gamesResponse []resources.Game
	w.Header().Set("Content-Type", "application/json")

	gamesModels, err := c.service.GetGames(*r.Date())
	if err != nil {
		c.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, game := range gamesModels {
		gamesResponse = append(gamesResponse, resources.NewGameResource(game))
	}
	if err := json.NewEncoder(w).Encode(gamesResponse); err != nil {
		c.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetGame
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *Controller) GetGame(w http.ResponseWriter, r *requests.GetSpecificGameRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetGame(r.Id())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := resources.NewGameResource(*gameModel)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetGameMetrics
//
// retrieve specific game and then calculate IMP metrics for every player
func (c *Controller) GetGameMetrics(w http.ResponseWriter, r *requests.GetSpecificGameMetricsRequest) {
	gameModel, err := c.service.GetGameMetrics(r.Id(), r.Pers())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	formatter := formatters.NewFormatter(r.Format())

	if err := formatter.Format(w, gameModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
