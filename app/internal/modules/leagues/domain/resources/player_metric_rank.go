package resources

import "IMP/app/internal/modules/players/domain/resources"

type PlayerMetricRank struct {
	Rank            int                 `json:"rank"`
	PlayerImpMetric resources.AvgMetric `json:"player"`
}
