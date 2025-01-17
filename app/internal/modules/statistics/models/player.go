package models

import "time"

type PlayerDTO struct {
	FullName       string
	BirthDate      *time.Time
	LeaguePlayerID int
	Statistic      PlayerStatisticDTO
}
