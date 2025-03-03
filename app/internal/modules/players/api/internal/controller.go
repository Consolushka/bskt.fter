package internal

import (
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
