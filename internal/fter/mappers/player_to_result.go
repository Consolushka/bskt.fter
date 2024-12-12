package mappers

import (
	"FTER/internal/fter/models"
	"FTER/internal/fter/results"
)

func PlayerToResult(player models.PlayerModel, fter float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player: player,
		FTER:   fter,
	}
}
