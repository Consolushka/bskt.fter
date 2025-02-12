package models

import (
	"IMP/app/internal/modules/leagues"
	"time"
)

type Team struct {
	ID         int            `json:"id" db:"id"`
	Alias      string         `json:"alias" db:"alias"`
	LeagueID   int            `json:"league_id" db:"league_id"`
	League     leagues.League `json:"league" gorm:"foreignKey:LeagueID"`
	Name       string         `json:"name" db:"name"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`
	OfficialId string         `json:"official_id" db:"official_id"`
}

func (Team) TableName() string {
	return "teams"
}
