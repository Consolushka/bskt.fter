package internal

import (
	"IMP/app/internal/abstract"
	gamesResources "IMP/app/internal/modules/games/domain/resources"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/modules/teams/api/internal/requests"
	teamsResponses "IMP/app/internal/modules/teams/api/responses"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Controller struct {
	abstract.BaseController

	service *teams.Service

	logger *logrus.Logger
}

func NewController() *Controller {
	return &Controller{
		service: teams.NewService(),
	}
}

// GetTeams returns all games filtered by date
func (c *Controller) GetTeams(w http.ResponseWriter, r *requests.GetTeamsRequest) {
	var response []teamsResponses.TeamResponse

	teamsArray, err := c.service.GetTeams()
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	for _, team := range teamsArray {
		response = append(response, teamsResponses.NewTeamResponse(team))
	}

	c.Ok(w, response)
}

// GetTeam
//
// id is the ID of the game to retrieve, could be only int. If not int, returns BadRequest
//
// id should be an id of existing game. If there is no game with given id, returns InternalServerError
func (c *Controller) GetTeam(w http.ResponseWriter, r *requests.GetTeamByIdRequest) {
	team, err := c.service.GetTeamById(r.Id())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	result := teamsResponses.NewTeamResponse(team)

	c.Ok(w, result)
}

// GetTeamGames
//
// Retrieve all games box score for given team
func (c *Controller) GetTeamGames(w http.ResponseWriter, r *requests.GetTeamGamesRequest) {
	var response []gamesResources.Game

	gamesModel, err := c.service.GetTeamWithGames(r.Id())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	for _, game := range gamesModel {
		response = append(response, gamesResources.NewGameResource(*game))
	}

	c.Ok(w, response)
}

func (c *Controller) GetTeamGamesMetrics(w http.ResponseWriter, r *requests.GetTeamByIdGamesMetricsRequest) {
	var response []gamesResources.Metric

	gameModel, err := c.service.GetTeamWithGamesMetrics(r.Id(), r.Pers())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	for _, game := range gameModel {
		response = append(response, gamesResources.NewMetricResource(*game))
	}

	c.Ok(w, response)
}
