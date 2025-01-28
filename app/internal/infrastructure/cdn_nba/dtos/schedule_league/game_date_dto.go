package schedule_league

type GameDateDTO struct {
	GameDate string    `json:"gameDate"`
	Games    []GameDTO `json:"games"`
}
