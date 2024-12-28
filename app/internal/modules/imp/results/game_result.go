package results

type TeamResults struct {
	Title string

	Players []PlayerImpResult
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
