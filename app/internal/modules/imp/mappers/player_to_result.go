package mappers

import (
	"FTER/app/internal/modules/imp/enums"
	"FTER/app/internal/modules/imp/models"
	"FTER/app/internal/modules/imp/results"
)

func PlayerToResult(player models.PlayerModel, timeBases []enums.TimeBasedImpCoefficient, impPersResults []float64) results.PlayerFterResult {
	return results.PlayerFterResult{
		Player:         player,
		TimeBases:      timeBases,
		ImpPersResults: impPersResults,
	}
}
