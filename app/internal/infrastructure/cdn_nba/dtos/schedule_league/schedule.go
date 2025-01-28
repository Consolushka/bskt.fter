package schedule_league

type SeasonScheduleDTO struct {
	SeasonYear string        `json:"seasonYear"`
	LeagueId   string        `json:"leagueId"`
	Games      []GameDateDTO `json:"gameDates"`
	Weeks      []WeekDTO     `json:"weeks"`
}
