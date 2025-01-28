package schedule_league

type WeekDTO struct {
	WeekNumber int    `json:"weekNumber"`
	WeekName   string `json:"weekName"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}
