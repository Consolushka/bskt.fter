package games

import (
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/teams"
	"time"
)

type GameModel struct {
	ID            int            `json:"id" gorm:"primaryKey"`
	LeagueID      int            `json:"league_id" gorm:"not null"`
	League        leagues.League `json:"league" gorm:"foreignKey:LeagueID"`
	HomeTeamID    int            `json:"home_team_id" gorm:"not null"`
	HomeTeam      teams.Team     `json:"home_team" gorm:"foreignKey:HomeTeamID"`
	AwayTeamID    int            `json:"away_team_id" gorm:"not null"`
	AwayTeam      teams.Team     `json:"away_team" gorm:"foreignKey:AwayTeamID"`
	PlayedMinutes int            `json:"played_minutes" gorm:"not null"`
	ScheduledAt   time.Time      `json:"scheduled_at" gorm:"not null"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (GameModel) TableName() string {
	return "games"
}
