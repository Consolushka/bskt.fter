package mappers

import (
	"IMP/app/internal/modules/imp/enums"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/imp/results"
)

func PlayerToResult(player models.PlayerModel, timeBases []enums.TimeBasedImpCoefficient, impPersResults []float64) results.PlayerImpResult {
	return results.PlayerImpResult{
		Player:         player,
		TimeBases:      timeBases,
		ImpPersResults: impPersResults,
	}
}
