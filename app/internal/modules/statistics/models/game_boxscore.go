package models

import (
	"time"
)

type GameBoxScoreDTO struct {
	Id            string
	LeagueAliasEn string
	HomeTeam      TeamBoxScoreDTO
	AwayTeam      TeamBoxScoreDTO
	PlayedMinutes int
	ScheduledAt   time.Time
}
