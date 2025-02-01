package games

import (
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/teams"
	"time"
)

type GameModel struct {
	ID            int                 `json:"id" gorm:"primaryKey"`
	LeagueID      int                 `json:"league_id" gorm:"not null"`
	League        leagues.League      `json:"league" gorm:"foreignKey:LeagueID"`
	HomeTeamID    int                 `json:"home_team_id" gorm:"not null"`
	HomeTeamStats teams.TeamGameStats `json:"home_team_stats" gorm:"foreignKey:GameId,TeamId;references:ID,HomeTeamID"`
	AwayTeamID    int                 `json:"away_team_id" gorm:"not null"`
	AwayTeamStats teams.TeamGameStats `json:"away_team_stats" gorm:"foreignKey:GameId,TeamId;references:ID,AwayTeamID"`
	PlayedMinutes int                 `json:"played_minutes" gorm:"not null"`
	ScheduledAt   time.Time           `json:"scheduled_at" gorm:"not null"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	OfficialId    int                 `json:"official_id" gorm:"not null"`
}

func (GameModel) TableName() string {
	return "games"
}
