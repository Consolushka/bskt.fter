package internal

import (
	"NBATrueEfficency/internal/statistics/repositories/sport_radar/dtos"
	"NBATrueEfficency/internal/true_efficiency/models"
)

func DTOToPlayerStats(dto dtos.PlayerStatsDTO) *models.PlayerStatsModel {
	return &models.PlayerStatsModel{
		Points:            dto.Points,
		OffensiveRtg:      dto.OffensiveRating,
		TwoPointsMade:     dto.TwoPointsMade,
		TwoPointsAtt:      dto.TwoPointsAtt,
		ThreePointsMade:   dto.ThreePointsMade,
		ThreePointsAtt:    dto.ThreePointsAtt,
		BlockedAtt:        dto.BlockedAtt,
		FreeThrowsMade:    dto.FreeThrowsMade,
		FreeThrowsAtt:     dto.FreeThrowsAtt,
		OffensiveRebounds: dto.OffensiveRebounds,
		DefensiveRebounds: dto.DefensiveRebounds,
		Assists:           dto.Assists,
		Turnovers:         dto.Turnovers,
		Steals:            dto.Steals,
		Blocks:            dto.Blocks,
		DefensiveRtg:      dto.DefensiveRating,
		PersonalFouls:     dto.PersonalFouls,
		TechFouls:         dto.TechFouls,
		FlagrantFouls:     dto.FlagrantFouls,
		PlsMin:            dto.PlsMin,
		FoulsDrawn:        dto.FoulsDrawn,
		OffensiveFouls:    dto.OffensiveFouls,
	}
}
