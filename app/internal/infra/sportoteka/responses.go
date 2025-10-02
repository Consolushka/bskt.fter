package sportoteka

type GameBoxScoreResponse struct {
	Status  string             `json:"status"`
	Message interface{}        `json:"message"`
	Result  GameBoxScoreEntity `json:"result"`
}

type CalendarResponse struct {
	Index      int                  `json:"index"`
	Status     string               `json:"status"`
	Message    interface{}          `json:"message"`
	TotalCount int                  `json:"totalCount"`
	Items      []CalendarGameEntity `json:"items"`
}
