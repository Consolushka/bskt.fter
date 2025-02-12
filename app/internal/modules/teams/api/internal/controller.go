package internal

import (
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/modules/teams/api/internal/requests"
	"IMP/app/log"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller struct {
	service *teams.Service

	logger *logrus.Logger
}

func NewController() *Controller {
	return &Controller{
		service: teams.NewService(),
		logger:  log.GetLogger(),
	}
}

// GetTeams returns all games filtered by date
func (c *Controller) GetTeams(w http.ResponseWriter, r *requests.GetTeamsRequest) {
	w.Header().Set("Content-Type", "application/json")

	games, err := c.service.GetTeams()
	if err != nil {
		c.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(games); err != nil {
		c.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTeam
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *Controller) GetTeam(w http.ResponseWriter, r *requests.GetTeamByIdRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetTeamById(r.Id())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(gameModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTeamGames
//
// retrieve specific game and then calculate IMP metrics for every player
func (c *Controller) GetTeamGames(w http.ResponseWriter, r *requests.GetTeamGamesRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetTeamWithGames(r.Id())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(gameModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetTeamGamesMetrics(w http.ResponseWriter, r *requests.GetTeamByIdGamesMetricsRequest) {
	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetTeamWithGamesMetrics(r.Id(), r.Pers())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(gameModel); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
