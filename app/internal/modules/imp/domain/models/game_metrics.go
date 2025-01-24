package models

import (
	"IMP/app/internal/modules/imp/domain/enums"
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
	FullNameLocal string
	FullNameEn    string
	SecondsPlayed int
	PlsMin        int
	ImpPers       []PlayerImpPersMetrics
}

type PlayerImpPersMetrics struct {
	Per enums.ImpPERs
	IMP float64
}

func (p *PlayerImpMetrics) Headers() []string {
	headers := []string{"Player",
		"Minutes Played",
	}

	for _, per := range p.ImpPers {
		headers = append(headers, string(per.Per))
	}

	return headers
}

func (p *PlayerImpMetrics) ToTable() []string {
	table := []string{
		p.FullNameEn,
		time_utils.SecondsToFormat(p.SecondsPlayed, time_utils.PlayedTimeFormat),
	}

	for _, per := range p.ImpPers {
		table = append(table, strconv.FormatFloat(per.IMP, 'f', 2, 64))
	}

	return table
}
