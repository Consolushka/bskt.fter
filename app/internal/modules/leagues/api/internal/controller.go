package internal

import (
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/leagues/api/internal/requests"
	"IMP/app/internal/modules/leagues/api/internal/responses"
	"IMP/app/internal/modules/leagues/domain/resources"
	"encoding/json"
	"net/http"
)

type Controller struct {
	service *leagues.Service
}

func NewController() *Controller {
	return &Controller{
		service: leagues.NewService(),
	}
}

func (c *Controller) GetLeagues(w http.ResponseWriter, r *requests.GetLeaguesRequest) {
	var response []resources.LeagueResource

	leagueModels, err := c.service.GetAllLeagues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, league := range leagueModels {
		response = append(response, resources.NewLeagueResponse(&league))
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) GetGamesByLeagueAndDate(w http.ResponseWriter, r *requests.GetGamesByLeagueAndDate) {
	gamesModel, err := c.service.GetGamesByLeagueAndDate(r.Id(), *r.Date())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewGamesByDateResponse(*r.Date(), gamesModel)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
