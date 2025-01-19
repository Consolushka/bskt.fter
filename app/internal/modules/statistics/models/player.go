package models

import "time"

type PlayerDTO struct {
	FullNameLocal  string
	FullNameEn     string
	BirthDate      *time.Time
	LeaguePlayerID int
	Statistic      PlayerStatisticDTO
}
