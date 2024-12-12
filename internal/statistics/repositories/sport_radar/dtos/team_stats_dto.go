package dtos

import (
	models2 "NBATrueEfficency/internal/fter/models"
	"NBATrueEfficency/internal/utils/arrays"
)

type TeamStatsDTO struct {
	Alias   string      `json:"alias"`
	Points  int         `json:"points"`
	Players []PlayerDTO `json:"players"`
}

func (dto TeamStatsDTO) ToFterModel() models2.TeamGameResultModel {
	return models2.TeamGameResultModel{
		Team: models2.TeamModel{
			FullName: dto.Alias,
			Alias:    dto.Alias,
		},
		TotalPoints: dto.Points,
		Players:     gameDtoToPlayers(dto.Players),
	}
}

func gameDtoToPlayers(players []PlayerDTO) []models2.PlayerModel {
	players = arrays.Filter(players, func(p PlayerDTO) bool {
		return p.Played
	})

	playersModels := make([]models2.PlayerModel, len(players))

	for i, player := range players {

		playersModels[i] = models2.PlayerModel{
			FullName:      player.FullName,
			MinutesPlayed: player.Statistics.Minutes,
			PlsMin:        player.Statistics.PlsMin,
		}
	}
	return playersModels
}
