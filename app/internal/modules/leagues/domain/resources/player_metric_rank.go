package resources

type PlayerMetricRank struct {
	Rank             int     `json:"rank"`
	FullName         string  `json:"full_name"`
	AvgMinutesPlayed string  `json:"avg_minutes_played"`
	GamesPlayed      int     `json:"games_played"`
	ImpPer           float64 `json:"imp_pers"`
}
