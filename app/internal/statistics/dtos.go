package statistics

import "time"

type GameBoxScoreDTO struct {
	Id            string
	LeagueAliasEn string
	IsFinal       bool
	HomeTeam      TeamBoxScoreDTO
	AwayTeam      TeamBoxScoreDTO
	PlayedMinutes int
	ScheduledAt   time.Time
}

type PlayerDTO struct {
	FullNameLocal  string
	FullNameEn     string
	BirthDate      *time.Time
	LeaguePlayerID string
	Statistic      PlayerStatisticDTO
}

type PlayerStatisticDTO struct {
	PlsMin        int
	PlayedSeconds int
	IsBench       bool
}

type TeamBoxScoreDTO struct {
	Alias    string
	Name     string
	LeagueId string
	Scored   int
	Players  []PlayerDTO
}
