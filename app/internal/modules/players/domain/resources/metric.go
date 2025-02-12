package resources

import (
	"IMP/app/internal/modules/imp/domain/models"
	impResources "IMP/app/internal/modules/imp/domain/resources"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
)

type Metric struct {
	FullName      string                `json:"full_name"`
	MinutesPlayed string                `json:"minutes_played"`
	PlsMin        int                   `json:"pls_min"`
	IMP           []impResources.Metric `json:"impPers"`
}

func NewMetricResource(player models.PlayerImpMetrics) Metric {
	return Metric{
		FullName:      player.FullNameLocal,
		MinutesPlayed: time_utils.SecondsToFormat(player.SecondsPlayed, time_utils.PlayedTimeFormat),
		PlsMin:        player.PlsMin,
		IMP: array_utils.Map(player.ImpPers, func(impPer models.PlayerImpPersMetrics) impResources.Metric {
			return impResources.Metric{
				Base: string(impPer.Per),
				Imp:  impPer.IMP,
			}
		}),
	}
}
