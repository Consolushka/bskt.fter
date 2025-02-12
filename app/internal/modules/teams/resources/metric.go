package resources

import (
	"IMP/app/internal/modules/imp/domain/models"
	"IMP/app/internal/modules/players/domain/resources"
	"IMP/app/internal/utils/array_utils"
)

type Metric struct {
	Alias       string             `json:"alias"`
	TotalPoints int                `json:"total_points"`
	Players     []resources.Metric `json:"players"`
}

func NewMetricResource(team models.TeamImpMetrics) Metric {
	return Metric{
		Alias:       team.Alias,
		TotalPoints: team.TotalPoints,
		Players: array_utils.Map(team.Players, func(player models.PlayerImpMetrics) resources.Metric {
			return resources.NewMetricResource(player)
		}),
	}
}
