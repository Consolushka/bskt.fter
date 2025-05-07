package infobasket

type GameBoxScoreResponse struct {
	IsOnline          bool              `json:"IsOnline"`
	GameStatus        int               `json:"GameStatus"`
	MaxPeriod         int               `json:"MaxPeriod"`
	FromDate          interface{}       `json:"FromDate"`
	GameDate          string            `json:"GameDate"`
	HasTime           bool              `json:"HasTime"`
	GameTime          string            `json:"GameTime"`
	GameTimeMsk       string            `json:"GameTimeMsk"`
	HasVideo          bool              `json:"HasVideo"`
	GameTeams         []TeamBoxScoreDto `json:"GameTeams"`
	CompNameRu        string            `json:"CompNameRu"`
	CompNameEn        string            `json:"CompNameEn"`
	LeagueNameRu      string            `json:"LeagueNameRu"`
	LeagueNameEn      string            `json:"LeagueNameEn"`
	LeagueShortNameRu string            `json:"LeagueShortNameRu"`
	LeagueShortNameEn string            `json:"LeagueShortNameEn"`
	Gender            int               `json:"Gender"`
	CompID            int               `json:"CompID"`
	LeagueID          int               `json:"LeagueID"`
	Is3x3             bool              `json:"Is3x3"`
}

type TeamScheduleResponse struct {
	Games []GameScheduleDto `json:"Games"`
}

type SeasonScheduleResponse struct {
	Games []GameScheduleDto `json:"Games"`
}
