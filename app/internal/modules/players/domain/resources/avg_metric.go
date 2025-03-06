package resources

import "IMP/app/internal/modules/imp/domain/resources"

type AvgMetric struct {
	FullName         string             `json:"full_name"`
	AvgMinutesPlayed string             `json:"avg_minutes_played"`
	GamesPlayed      int                `json:"games_played"`
	ImpPers          []resources.Metric `json:"imp_pers"`
}
