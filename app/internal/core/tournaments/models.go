package tournaments

import (
	"IMP/app/internal/core/leagues"
	"time"

	"gorm.io/gorm"
)

type TournamentModel struct {
	Id                 uint           `gorm:"column:id"`
	LeagueId           uint           `gorm:"column:league_id"`
	Name               string         `gorm:"column:name"`
	StartAt            time.Time      `gorm:"column:start_at"`
	EndAt              time.Time      `gorm:"column:end_at"`
	RegulationDuration int            `gorm:"column:regulation_duration"`
	CreatedAt          time.Time      `gorm:"column:created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at"`

	League   leagues.LeagueModel `gorm:"foreignKey:LeagueId"`
	Provider TournamentProvider  `gorm:"foreignKey:TournamentId"`
}

func (TournamentModel) TableName() string {
	return "tournaments"
}

type TournamentProvider struct {
	TournamentId uint      `gorm:"column:tournament_id"`
	ProviderName string    `gorm:"column:provider_name"`
	ExternalId   *string   `gorm:"column:external_id"`
	Params       []byte    `gorm:"column:params;type:json"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (TournamentProvider) TableName() string {
	return "tournament_providers"
}
