package resources

import (
	"IMP/app/internal/modules/imp/domain/models"
	"IMP/app/internal/modules/teams/resources"
	"strconv"
)

type Metric struct {
	GameId       string           `json:"game_id"`
	Scheduled    string           `json:"scheduled"`
	FullGameTime string           `json:"full_game_time"`
	Home         resources.Metric `json:"home"`
	Away         resources.Metric `json:"away"`
}

func NewMetricResource(metrics models.GameImpMetrics) Metric {
	return Metric{
		GameId:       strconv.Itoa(metrics.Id),
		Scheduled:    metrics.Scheduled.Format("02.01.2006 15:04"),
		FullGameTime: strconv.Itoa(metrics.FullGameTime),
		Home:         resources.NewMetricResource(metrics.Home),
		Away:         resources.NewMetricResource(metrics.Away),
	}
}
