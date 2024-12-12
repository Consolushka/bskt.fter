package mappers

import (
	"NBATrueEfficency/internal/fter/models"
	"NBATrueEfficency/internal/fter/results"
)

func PlayerToResult(player models.PlayerModel, fter float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player: player,
		FTER:   fter,
	}
}
