package api_nba

type GamesResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		Date string `json:"date"`
	} `json:"parameters"`
	Errors   []interface{} `json:"errors"`
	Results  int           `json:"results"`
	Response []GameEntity  `json:"response"`
}

type PlayerStatisticResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		Game string `json:"game"`
	} `json:"parameters"`
	Errors   []interface{}           `json:"errors"`
	Results  int                     `json:"results"`
	Response []PlayerStatisticEntity `json:"response"`
}

type PlayersResponse struct {
	Get        string `json:"get"`
	Parameters struct {
		Id string `json:"id"`
	} `json:"parameters"`
	Errors   []interface{}  `json:"errors"`
	Results  int            `json:"results"`
	Response []PlayerEntity `json:"response"`
}
