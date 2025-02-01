package models

type TeamBoxScoreDTO struct {
	Alias    string
	Name     string
	LeagueId string
	Scored   int
	Players  []PlayerDTO
}
