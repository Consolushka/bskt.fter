package models

import (
	"time"
)

type GameImpMetrics struct {
	Id           int
	Scheduled    *time.Time
	Home         TeamImpMetrics
	Away         TeamImpMetrics
	FullGameTime int
}

type TeamImpMetrics struct {
	Alias       string
	TotalPoints int
	Players     []PlayerImpMetrics
}

type PlayerImpMetrics struct {
	FullName      string
	SecondsPlayed int
	PlsMin        int
	IMP           float64
}
