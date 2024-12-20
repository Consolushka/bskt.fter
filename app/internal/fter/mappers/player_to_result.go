package mappers

import (
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
)

func PlayerToResult(player models.PlayerModel, fter float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player: player,
		FTER:   fter,
	}
}
