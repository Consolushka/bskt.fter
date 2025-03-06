package responses

import "IMP/app/internal/modules/leagues/domain/resources"

type RankingResponse struct {
	Ranking []resources.PlayerMetricRank `json:"ranking"`
}
