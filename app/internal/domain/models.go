package domain

import (
	"time"
)

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

type Team struct {
	ID         int       `json:"id" db:"id"`
	Alias      string    `json:"alias" db:"alias"`
	LeagueID   int       `json:"league_id" db:"league_id"`
	League     League    `json:"league" gorm:"foreignKey:LeagueID"`
	Name       string    `json:"name" db:"name"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	OfficialId string    `json:"official_id" db:"official_id"`
}

type TeamGameStats struct {
	Id              int               `json:"id" db:"id"`
	TeamId          int               `json:"team_id" db:"team_id"`
	Team            *Team             `json:"team" gorm:"foreignKey:TeamId"`
	GameId          int               `json:"game_id" db:"game_id"`
	Points          int               `json:"points" db:"points"`
	PlayerGameStats []PlayerGameStats `json:"players" gorm:"foreignKey:team_game_id"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

type Game struct {
	ID            int           `json:"id" gorm:"primaryKey"`
	LeagueID      int           `json:"league_id" gorm:"not null"`
	League        League        `json:"league" gorm:"foreignKey:LeagueID"`
	HomeTeamID    int           `json:"home_team_id" gorm:"not null"`
	HomeTeamStats TeamGameStats `json:"home_team_stats" gorm:"foreignKey:GameId,TeamId;references:ID,HomeTeamID"`
	AwayTeamID    int           `json:"away_team_id" gorm:"not null"`
	AwayTeamStats TeamGameStats `json:"away_team_stats" gorm:"foreignKey:GameId,TeamId;references:ID,AwayTeamID"`
	PlayedMinutes int           `json:"played_minutes" gorm:"not null"`
	ScheduledAt   time.Time     `json:"scheduled_at" gorm:"not null"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	OfficialId    string        `json:"official_id" gorm:"not null"`
}

type Player struct {
	ID            int        `json:"id" db:"id"`
	FullNameLocal string     `json:"full_name" db:"full_name_local"`
	FullNameEn    string     `json:"full_name_eng" db:"full_name_en"`
	BirthDate     *time.Time `json:"birth_date" db:"birth_date"`
	OfficialId    string     `json:"official_id" db:"official_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type PlayerGameStats struct {
	PlayerID      int       `json:"player_id" db:"player_id" gorm:"primaryKey;autoIncrement:false"`
	Player        *Player   `json:"player" gorm:"foreignKey:PlayerID"`
	TeamGameId    int       `json:"team_game_id" db:"team_game_id" gorm:"index"`
	PlsMin        int       `json:"pls_min" db:"pls_min"`
	PlayedSeconds int       `json:"played_min" db:"played_seconds"`
	IsBench       bool      `json:"is_bench" db:"is_bench"`
	IMPClean      float64   `json:"imp_clean" db:"imp_clean"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
