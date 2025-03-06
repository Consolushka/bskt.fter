package internal

import (
	"IMP/app/internal/base/components/request_components"
	"IMP/app/internal/modules/players"
	"IMP/app/internal/modules/players/api/internal/requests"
	"IMP/app/internal/modules/players/api/internal/responses"
	"IMP/app/internal/modules/players/domain/models"
	"IMP/app/internal/utils/array_utils"
	"encoding/json"
	"net/http"
)

type Controller struct {
	service *players.Service
}

func NewController() *Controller {
	return &Controller{
		service: players.NewService(),
	}
}

// Search players by local or en full name
func (c *Controller) Search(w http.ResponseWriter, r *requests.SearchPlayerByFullNameRequest) {
	var response []responses.FoundPlayerResponse

	playersModels, err := c.service.GetPlayerByFullName(r.FullName())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response = array_utils.Map(playersModels, func(playerModel models.Player) responses.FoundPlayerResponse {
		return responses.FoundPlayerResponse{
			Id:            playerModel.ID,
			FullNameEn:    playerModel.FullNameEn,
			FullNameLocal: playerModel.FullNameLocal,
			BirthDate:     playerModel.BirthDate.Format("02-01-2006"),
		}
	})

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) PlayerGamesBoxScore(w http.ResponseWriter, r *request_components.HasIdPathParam) {
	games, err := c.service.GetPlayerGamesBoxScore(r.Id())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewPlayerGamesBoxScoreResponse(games)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *Controller) PlayerGamesMetrics(w http.ResponseWriter, r *requests.PlayerGamesMetricsRequest) {
	games, err := c.service.GetPlayerGamesMetrics(r.Id(), r.Pers())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responses.NewPlayerGamesMetricsResponse(games)

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
