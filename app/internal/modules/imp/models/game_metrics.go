package models

import (
	"IMP/app/internal/utils/time_utils"
	"strconv"
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

func (p *PlayerImpMetrics) Headers() []string {
	return []string{"Player",
		"Minutes Played",
		"IMP",
	}
}

func (p *PlayerImpMetrics) ToTable() []string {
	return []string{
		p.FullName,
		time_utils.SecondsToFormat(p.SecondsPlayed, time_utils.PlayedTimeFormat),
		strconv.FormatFloat(p.IMP, 'f', 2, 64),
	}
}
