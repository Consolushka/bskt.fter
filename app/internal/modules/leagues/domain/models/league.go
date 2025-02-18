package models

import "time"

//todo: add constraints for leagues

const MLBLAlias = "mlbl"
const NBAAlias = "nba"

type League struct {
	ID               int       `db:"id" json:"id"`
	NameLocal        string    `db:"name_local" json:"name_local"`
	AliasLocal       string    `db:"alias_local" json:"alias_local"`
	NameEn           string    `db:"name_en" json:"name_en"`
	AliasEn          string    `db:"alias_en" json:"alias_en"`
	PeriodsNumber    int       `db:"periods_number" json:"periods_number"`
	PeriodDuration   int       `db:"period_duration" json:"period_duration"`
	OvertimeDuration int       `db:"overtime_duration" json:"overtime_duration"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}
