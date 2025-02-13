package resources

import (
	"IMP/app/internal/modules/imp/domain/models"
	teamsResources "IMP/app/internal/modules/teams/domain/resources"
	"strconv"
)

type Metric struct {
	GameId       string                `json:"game_id"`
	Scheduled    string                `json:"scheduled"`
	FullGameTime string                `json:"full_game_time"`
	Home         teamsResources.Metric `json:"home"`
	Away         teamsResources.Metric `json:"away"`
}

func NewMetricResource(metrics models.GameImpMetrics) Metric {
	return Metric{
		GameId:       strconv.Itoa(metrics.Id),
		Scheduled:    metrics.Scheduled.Format("02.01.2006 15:04"),
		FullGameTime: strconv.Itoa(metrics.FullGameTime),
		Home:         teamsResources.NewMetricResource(metrics.Home),
		Away:         teamsResources.NewMetricResource(metrics.Away),
	}
}
