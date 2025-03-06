package internal

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/modules/games"
	"IMP/app/internal/modules/games/api/internal/formatters"
	"IMP/app/internal/modules/games/api/internal/requests"
	"IMP/app/internal/modules/games/domain/resources"
	"net/http"
)

type Controller struct {
	abstract.BaseController
	service *games.Service
}

func NewController() *Controller {
	return &Controller{
		service: games.NewService(),
	}
}

// GetGames returns all games filtered by date
func (c *Controller) GetGames(w http.ResponseWriter, r *requests.GetGamesRequest) {
	var gamesResponse []resources.Game

	gamesModels, err := c.service.GetGames(*r.Date())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	for _, game := range gamesModels {
		gamesResponse = append(gamesResponse, resources.NewGameResource(game))
	}

	c.Ok(w, gamesResponse)
}

// GetGame
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *Controller) GetGame(w http.ResponseWriter, r *requests.GetSpecificGameRequest) {
	gameModel, err := c.service.GetGame(r.Id())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	response := resources.NewGameResource(*gameModel)

	c.Ok(w, response)
}

// GetGameMetrics
//
// retrieve specific game and then calculate IMP metrics for every player
func (c *Controller) GetGameMetrics(w http.ResponseWriter, r *requests.GetSpecificGameMetricsRequest) {
	gameModel, err := c.service.GetGameMetrics(r.Id(), r.Pers())

	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	formatter := formatters.NewFormatter(r.Format())

	if err := formatter.Format(w, gameModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
