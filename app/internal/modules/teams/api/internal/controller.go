package internal

import (
	"IMP/app/internal/modules/games/resources"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/modules/teams/api/internal/requests"
	teamResponses "IMP/app/internal/modules/teams/api/responses"
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
	}
}

// GetTeams returns all games filtered by date
func (c *Controller) GetTeams(w http.ResponseWriter) {
	var response []teamResponses.TeamResponse

	w.Header().Set("Content-Type", "application/json")

	teamsArray, err := c.service.GetTeams()
	if err != nil {
		c.logger.Errorln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, team := range teamsArray {
		response = append(response, teamResponses.NewTeamResponse(team))
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
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

	team, err := c.service.GetTeamById(r.Id())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := teamResponses.NewTeamResponse(team)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetTeamGames
//
// retrieve specific game and then calculate IMP metrics for every player
func (c *Controller) GetTeamGames(w http.ResponseWriter, r *requests.GetTeamGamesRequest) {
	var response []resources.Game

	w.Header().Set("Content-Type", "application/json")

	gamesModel, err := c.service.GetTeamWithGames(r.Id())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, game := range gamesModel {
		response = append(response, resources.NewGameResource(*game))
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetTeamGamesMetrics(w http.ResponseWriter, r *requests.GetTeamByIdGamesMetricsRequest) {
	var response []resources.Metric

	w.Header().Set("Content-Type", "application/json")

	gameModel, err := c.service.GetTeamWithGamesMetrics(r.Id(), r.Pers())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, game := range gameModel {
		response = append(response, resources.NewMetricResource(*game))
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
