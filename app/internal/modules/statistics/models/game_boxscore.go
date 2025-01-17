package models

import (
	"IMP/app/internal/modules/statistics/enums"
	"time"
)

type GameBoxScoreDTO struct {
	League        enums.League
	HomeTeam      TeamBoxScoreDTO
	AwayTeam      TeamBoxScoreDTO
	PlayedMinutes int
	ScheduledAt   time.Time
}
