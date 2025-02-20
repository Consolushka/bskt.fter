package cdn_nba

type GameBoxScoreResponse struct {
	Meta MetaDto     `json:"meta"`
	Game BoxScoreDto `json:"game"`
}

type ScheduleResponse struct {
	Meta     MetaDto           `json:"meta"`
	Schedule SeasonScheduleDto `json:"leagueSchedule"`
}
