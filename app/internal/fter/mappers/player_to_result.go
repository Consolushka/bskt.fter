package mappers

import (
	"FTER/app/internal/fter/enums"
	"FTER/app/internal/fter/results"
	"FTER/app/internal/models"
)

func PlayerToResult(player models.PlayerModel, timeBases []enums.TimeBasedImpCoefficient, impPersResults []float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player:         player,
		TimeBases:      timeBases,
		ImpPersResults: impPersResults,
	}
}
