package mappers

import (
	"FTER/internal/fter/results"
	"FTER/internal/models"
)

func PlayerToResult(player models.PlayerModel, fter float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player: player,
		FTER:   fter,
	}
}
