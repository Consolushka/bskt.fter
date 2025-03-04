package responses

import (
	gamesResources "IMP/app/internal/modules/games/domain/resources"
	"IMP/app/internal/modules/imp/domain/models"
	"IMP/app/internal/utils/array_utils"
)

type PlayerGamesMetricsResponse struct {
	Games []gamesResources.Metric `json:"games"`
}

func NewPlayerGamesMetricsResponse(games []*models.GameImpMetrics) PlayerGamesMetricsResponse {
	return PlayerGamesMetricsResponse{
		Games: array_utils.Map(games, func(game *models.GameImpMetrics) gamesResources.Metric {
			return gamesResources.NewMetricResource(*game)
		}),
	}

}
