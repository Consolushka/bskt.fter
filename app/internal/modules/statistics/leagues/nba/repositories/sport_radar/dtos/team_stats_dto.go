package dtos

import (
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/utils/array_utils"
)

type TeamStatsDTO struct {
	Alias   string      `json:"alias"`
	Points  int         `json:"points"`
	Players []PlayerDTO `json:"players"`
}

func (dto TeamStatsDTO) ToFterModel() models.TeamGameResultModel {
	return models.TeamGameResultModel{
		Team: models.TeamModel{
			FullName: dto.Alias,
			Alias:    dto.Alias,
		},
		TotalPoints: dto.Points,
		Players:     gameDtoToPlayers(dto.Players),
	}
}

func gameDtoToPlayers(players []PlayerDTO) []models.PlayerModel {
	players = array_utils.Filter(players, func(p PlayerDTO) bool {
		return p.Played
	})

	playersModels := make([]models.PlayerModel, len(players))

	for i, player := range players {
		playersModels[i] = player.ToFterModel()
	}

	return playersModels
}
