package dtos

import (
	"FTER/internal/models"
	"FTER/internal/utils/arrays"
)

type TeamDTO struct {
	TeamId            int         `json:"teamId"`
	TeamName          string      `json:"teamName"`
	TeamCity          string      `json:"teamCity"`
	TeamTricode       string      `json:"teamTricode"`
	Score             int         `json:"score"`
	InBonus           string      `json:"inBonus"`
	TimeoutsRemaining int         `json:"timeoutsRemaining"`
	Players           []PlayerDTO `json:"players"`
}

func (dto TeamDTO) ToFterModel() models.TeamGameResultModel {
	return models.TeamGameResultModel{
		Team: models.TeamModel{
			FullName: dto.TeamTricode,
			Alias:    dto.TeamTricode,
		},
		TotalPoints: dto.Score,
		Players:     gameDtoToPlayers(dto.Players),
	}
}

func gameDtoToPlayers(players []PlayerDTO) []models.PlayerModel {
	players = arrays.Filter(players, func(p PlayerDTO) bool {
		return p.Played == "1"
	})

	playersModels := make([]models.PlayerModel, len(players))

	for i, player := range players {
		playersModels[i] = player.ToFterModel()
	}

	return playersModels
}
