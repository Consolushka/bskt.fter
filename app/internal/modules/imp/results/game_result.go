package results

type TeamResults struct {
	Title string

	Players []PlayerFterResult
}

type GameResult struct {
	GameId   string
	Title    string
	Schedule string

	Home *TeamResults
	Away *TeamResults
}

func (r *GameResult) GetFileName() string {
	return r.Title + ". " + r.Schedule
}
