package api_basketball

type GamesResponse struct {
	Get        string       `json:"get"`
	Parameters interface{}  `json:"parameters"`
	Errors     []string     `json:"errors"`
	Results    int          `json:"results"`
	Response   []GameEntity `json:"response"`
}

type TeamStatsResponse struct {
	Get        string            `json:"get"`
	Parameters interface{}       `json:"parameters"`
	Errors     []string          `json:"errors"`
	Results    int               `json:"results"`
	Response   []TeamStatsEntity `json:"response"`
}

type PlayerStatsResponse struct {
	Get        string              `json:"get"`
	Parameters interface{}         `json:"parameters"`
	Errors     []string            `json:"errors"`
	Results    int                 `json:"results"`
	Response   []PlayerStatsEntity `json:"response"`
}

type PlayerInfoResponse struct {
	Get        string             `json:"get"`
	Parameters interface{}        `json:"parameters"`
	Errors     []string           `json:"errors"`
	Results    int                `json:"results"`
	Response   []PlayerInfoEntity `json:"response"`
}
