package models

type TeamBoxScoreDTO struct {
	Alias    string
	Name     string
	LeagueId int
	Players  []PlayerDTO
}
