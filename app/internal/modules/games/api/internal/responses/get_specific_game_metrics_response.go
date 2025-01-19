package responses

type GetSpecificGameMetricsResponse struct {
	GameId       string                             `json:"game_id"`
	Scheduled    string                             `json:"scheduled"`
	FullGameTime string                             `json:"full_game_time"`
	Home         GetSpecificGameTeamMetricsResponse `json:"home"`
	Away         GetSpecificGameTeamMetricsResponse `json:"away"`
}

type GetSpecificGameTeamMetricsResponse struct {
	Alias       string                                 `json:"alias"`
	TotalPoints int                                    `json:"total_points"`
	Players     []GetSpecificGamePlayerMetricsResponse `json:"players"`
}

type GetSpecificGamePlayerMetricsResponse struct {
	FullName      string  `json:"full_name"`
	MinutesPlayed string  `json:"minutes_played"`
	PlsMin        int     `json:"pls_min"`
	IMP           float64 `json:"imp"`
}
