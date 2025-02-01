package models

import (
	"IMP/app/internal/modules/statistics/enums"
	"time"
)

type GameBoxScoreDTO struct {
	Id            string
	League        enums.League
	HomeTeam      TeamBoxScoreDTO
	AwayTeam      TeamBoxScoreDTO
	PlayedMinutes int
	ScheduledAt   time.Time
}
