package todays_games

type ScoreboardDTO struct {
	GameDate   string    `json:"gameDate"`
	LeagueId   string    `json:"leagueId"`
	LeagueName string    `json:"leagueName"`
	Games      []GameDTO `json:"games"`
}
